import { getMenuButtonSafeArea } from '../../../utils/navigation';
import { PUBLISH_CATEGORIES } from '../data';
import { publishResource } from '../../../api/modules/resource';
import { getUploadResultPath, uploadFile } from '../../../api/modules/upload';

const MAX_FILES = 5;

function formatFileSize(bytes) {
  if (!bytes) return '0 B';
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
}

function getFileType(filePath) {
  if (!filePath) return 'PDF';
  const ext = filePath.split('.').pop().toLowerCase();
  const typeMap = {
    pdf: 'PDF',
    doc: 'Word',
    docx: 'Word',
    xls: 'Excel',
    xlsx: 'Excel',
    ppt: 'PPT',
    pptx: 'PPT',
    jpg: '图片',
    jpeg: '图片',
    png: '图片',
    gif: '图片',
    bmp: '图片',
    webp: '图片',
  };
  return typeMap[ext] || ext.toUpperCase();
}

function isImageFile(filePath) {
  if (!filePath) return false;
  const ext = filePath.split('.').pop().toLowerCase();
  return ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'].includes(ext);
}

function stripExtension(fileName) {
  if (!fileName) return '';
  const dotIndex = fileName.lastIndexOf('.');
  return dotIndex > 0 ? fileName.substring(0, dotIndex) : fileName;
}

Page({
  data: {
    canPublish: false,
    menuTop: 40,
    menuButtonHeight: 32,
    navHeight: 88,
    categoryList: PUBLISH_CATEGORIES,
    fileList: [],
    form: {
      title: '',
      category: '',
    },
    publishing: false,
  },

  onLoad() {
    this.applyNavigationSafeArea();
    this.updatePublishState();
  },

  applyNavigationSafeArea() {
    const { top, height } = getMenuButtonSafeArea(10);
    this.setData({
      menuTop: top,
      menuButtonHeight: height,
      navHeight: top + height + 16,
    });
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.redirectTo({ url: '/pages/resource/index' });
  },

  getEventValue(e) {
    if (e.detail && typeof e.detail === 'object' && 'value' in e.detail) return e.detail.value;
    return e.detail;
  },

  onFieldChange(e) {
    const { field } = e.currentTarget.dataset;
    if (!field) return;

    this.setData(
      {
        [`form.${field}`]: this.getEventValue(e),
      },
      () => {
        this.updatePublishState();
      },
    );
  },

  onCategorySelect(e) {
    const { value } = e.currentTarget.dataset;
    if (!value) return;

    this.setData({ 'form.category': value }, () => {
      this.updatePublishState();
    });
  },

  tryAutoFillTitle(newFiles) {
    if (this.data.form.title.trim()) return;

    const firstName = newFiles[0] && newFiles[0].name;
    if (!firstName) return;

    this.setData({ 'form.title': stripExtension(firstName) }, () => {
      this.updatePublishState();
    });
  },

  updatePublishState() {
    const { form, fileList } = this.data;
    const hasTitle = form.title.trim().length >= 1;
    const hasCategory = form.category.trim();
    const hasFiles = fileList.length > 0;

    this.setData({ canPublish: Boolean(hasTitle && hasCategory && hasFiles) });
  },

  addFilesToList(newFiles) {
    const merged = this.data.fileList.concat(newFiles).slice(0, MAX_FILES);
    this.setData({ fileList: merged }, () => {
      this.tryAutoFillTitle(newFiles);
      this.updatePublishState();
    });
  },

  chooseFromChat() {
    const remain = Math.max(0, MAX_FILES - this.data.fileList.length);
    if (!remain) return;

    wx.chooseMessageFile({
      count: remain,
      type: 'all',
      success: (res) => {
        const nextFiles = res.tempFiles
          .map((item) => {
            const filePath = item.path || item.tempFilePath || '';
            return {
              url: filePath,
              name: item.name || filePath.split('/').pop() || '未命名文件',
              size: formatFileSize(item.size),
              fileType: getFileType(item.name || filePath),
              isImage: isImageFile(item.name || filePath),
            };
          })
          .filter((item) => item.url);
        this.addFilesToList(nextFiles);
      },
    });
  },

  previewImage(e) {
    const index = Number(e.currentTarget.dataset.index || 0);
    const imageUrls = this.data.fileList.filter((f) => f.isImage).map((f) => f.url);
    if (!imageUrls.length) return;

    const currentUrl = this.data.fileList[index] && this.data.fileList[index].url;
    wx.previewImage({
      current: currentUrl,
      urls: imageUrls,
    });
  },

  removeFile(e) {
    const index = Number(e.currentTarget.dataset.index || 0);
    const fileList = this.data.fileList.filter((_, fileIndex) => fileIndex !== index);
    this.setData({ fileList }, () => {
      this.updatePublishState();
    });
  },

  async uploadFiles() {
    const { fileList } = this.data;
    if (fileList.some((file) => !file.url)) {
      throw new Error('文件路径无效，请重新选择文件');
    }

    let results = [];
    try {
      results = await Promise.all(
        fileList.map((file) => uploadFile(file.url, { name: file.name, scene: 'resource' })),
      );
    } catch (error) {
      if (error && error.statusCode === 404) {
        throw new Error('当前服务端未部署文件上传接口，暂时无法发布资料');
      }
      throw new Error(error.message || '文件上传失败，请稍后重试');
    }

    const filePaths = results.map(getUploadResultPath).filter(Boolean);

    if (filePaths.length !== fileList.length) {
      throw new Error('部分文件上传失败');
    }

    return filePaths;
  },

  async release() {
    const { form, fileList, canPublish, publishing } = this.data;

    if (publishing) return;

    if (!canPublish) {
      const missing = [];
      if (!form.title.trim()) missing.push('资料名');
      if (!form.category) missing.push('资料类型');
      if (!fileList.length) missing.push('上传文件');

      wx.showToast({ title: `请补充${missing.join('、')}`, icon: 'none' });
      return;
    }

    this.setData({ publishing: true });

    try {
      wx.showLoading({ title: '上传中' });
      const filePaths = await this.uploadFiles();

      const res = await publishResource({
        title: form.title.trim(),
        desc: form.title.trim(),
        category: form.category,
        course_name: '',
        contact: '',
        file_paths: filePaths,
      });
      const data = res.data || res || {};
      const insertId = data.id || '';

      wx.hideLoading();
      wx.showToast({ title: '发布成功', icon: 'success' });

      setTimeout(() => {
        if (getCurrentPages().length > 1) {
          wx.navigateBack({ delta: 1 });
          return;
        }
        wx.redirectTo({ url: `/pages/resource/index${insertId ? `?insertId=${insertId}` : ''}` });
      }, 360);
    } catch (e) {
      wx.hideLoading();
      const message = e.message || '发布失败，请重试';
      wx.showToast({ title: message, icon: 'none' });
    } finally {
      this.setData({ publishing: false });
    }
  },
});
