import { RESOURCE_CATEGORIES } from './data';
import { fetchResourceDetail, fetchResourceList, favoriteResource } from '../../api/modules/resource';

const CATEGORY_MAP = {};
RESOURCE_CATEGORIES.forEach((c) => {
  CATEGORY_MAP[c.value] = c.label;
});

const AVATAR_COLORS = ['blue', 'green', 'orange', 'purple', 'teal', 'rose'];

function stableIndex(value, modulo) {
  const text = `${value || ''}`;
  let total = 0;
  for (let i = 0; i < text.length; i += 1) total += text.charCodeAt(i);
  return total % modulo;
}

function getFileIcon(fileType) {
  const iconMap = {
    PDF: 'file-pdf',
    Word: 'file-word',
    Excel: 'file-excel',
    PPT: 'file-ppt',
    图片: 'image',
  };
  return iconMap[fileType] || 'file';
}

function getFileDisplayType(file = {}) {
  const rawType = file.file_type || '';
  const name = file.name || '';
  const ext = name.includes('.') ? name.split('.').pop().toLowerCase() : '';
  if (rawType.includes('pdf') || ext === 'pdf') return 'PDF';
  if (rawType.includes('word') || ['doc', 'docx'].includes(ext)) return 'Word';
  if (rawType.includes('excel') || rawType.includes('spreadsheet') || ['xls', 'xlsx'].includes(ext)) return 'Excel';
  if (rawType.includes('presentation') || ['ppt', 'pptx'].includes(ext)) return 'PPT';
  if (rawType.includes('image') || ['jpg', 'jpeg', 'png', 'gif', 'bmp', 'webp'].includes(ext)) return '图片';
  return rawType || (ext ? ext.toUpperCase() : '文件');
}

function mapResourceItem(item) {
  const extra = item.extra || {};
  const files = extra.files || [];
  const firstFile = files[0] || {};
  const category = extra.category || '';
  const fileType = getFileDisplayType(firstFile.file_type ? firstFile : { file_type: extra.file_type });
  return {
    id: item.id,
    title: item.title || '',
    desc: item.desc || '',
    category,
    categoryLabel: CATEGORY_MAP[category] || category || '其他',
    fileType,
    fileSize: firstFile.file_size || extra.file_size || '',
    files,
    downloadUrl: firstFile.url || extra.download_url || '',
    updatedAt: item.created_at || '',
    sponsor: item.publisher || '',
    avatarColor: AVATAR_COLORS[stableIndex(item.id, AVATAR_COLORS.length)],
    sponsorTag: extra.course_name || item.created_at || '',
    sponsorInitial: item.publisher_initial || (item.publisher ? item.publisher.charAt(0) : ''),
    views: extra.views || 0,
    likes: extra.likes || 0,
    fileIcon: getFileIcon(fileType),
  };
}

function moveItemToTop(items, id) {
  if (!id) return items;
  const target = items.find((item) => item.id === id);
  if (!target) return items;
  return [target].concat(items.filter((item) => item.id !== id));
}

function dedupeById(items) {
  const seen = {};
  return items.filter((item) => {
    if (!item.id || seen[item.id]) return false;
    seen[item.id] = true;
    return true;
  });
}

Page({
  data: {
    activeCategory: 'all',
    categoryList: RESOURCE_CATEGORIES,
    resourceList: [],
    visibleResources: [],
    searchKeyword: '',
    headerHeight: 194,
    page: 1,
    pageSize: 20,
    total: 0,
    loading: false,
    insertId: '',
  },

  onLoad(options = {}) {
    if (options.insertId) {
      this.setData({ insertId: options.insertId });
    }
  },

  onShow() {
    this.refreshResources();
  },

  async forceInsertResource(items) {
    const { insertId } = this.data;
    if (!insertId) return items;

    const moved = moveItemToTop(items, insertId);
    if (moved[0] && moved[0].id === insertId) return moved;

    try {
      const res = await fetchResourceDetail(insertId);
      const detail = mapResourceItem(res.data || res);
      return [detail].concat(items.filter((item) => item.id !== insertId));
    } catch (e) {
      return items;
    }
  },

  async refreshResources() {
    this.setData({ loading: true, page: 1 });
    const { activeCategory, searchKeyword, pageSize } = this.data;
    const params = { page: 1, pageSize };
    if (activeCategory !== 'all') params.category = activeCategory;
    if (searchKeyword && searchKeyword.trim()) params.keyword = searchKeyword.trim();

    try {
      const res = await fetchResourceList(params);
      const list = (res.data && res.data.list) || [];
      const total = (res.data && res.data.total) || 0;
      const mapped = await this.forceInsertResource(list.map(mapResourceItem));
      this.setData({
        resourceList: mapped,
        visibleResources: mapped,
        total: Math.max(total, mapped.length),
        page: 1,
        loading: false,
      });
    } catch (e) {
      this.setData({ loading: false });
    }
  },

  filterResources() {
    this.refreshResources();
  },

  onBack() {
    if (getCurrentPages().length > 1) {
      wx.navigateBack({ delta: 1 });
      return;
    }
    wx.switchTab({ url: '/pages/home/index' });
  },

  onSearchInput(e) {
    this.setData({ searchKeyword: e.detail.value }, () => {
      this.filterResources();
    });
  },

  onSearchConfirm() {
    this.filterResources();
  },

  onSearchClear() {
    this.setData({ searchKeyword: '' }, () => {
      this.filterResources();
    });
  },

  onCategoryChange(e) {
    const { value } = e.detail;
    if (!value || value === this.data.activeCategory) return;

    this.setData({ activeCategory: value }, () => {
      this.filterResources();
    });
  },

  onHeaderHeightChange(e) {
    const { height } = e.detail;
    if (height) this.setData({ headerHeight: height });
  },

  async onLoadMore() {
    const { loading, page, pageSize, total, activeCategory, searchKeyword } = this.data;
    if (loading) return;
    if ((page * pageSize) >= total) return;

    const nextPage = page + 1;
    this.setData({ loading: true });
    const params = { page: nextPage, pageSize };
    if (activeCategory !== 'all') params.category = activeCategory;
    if (searchKeyword && searchKeyword.trim()) params.keyword = searchKeyword.trim();

    try {
      const res = await fetchResourceList(params);
      const list = (res.data && res.data.list) || [];
      const mapped = list.map(mapResourceItem);
      const newVisible = dedupeById(this.data.visibleResources.concat(mapped));
      this.setData({
        resourceList: dedupeById(this.data.resourceList.concat(mapped)),
        visibleResources: newVisible,
        page: nextPage,
        loading: false,
      });
    } catch (e) {
      this.setData({ loading: false });
    }
  },

  goRelease() {
    wx.navigateTo({ url: '/pages/resource/publish/index' });
  },

  onOpenResource(e) {
    const resource = e.currentTarget.dataset.resource || {};
    const files = resource.files || [];
    const firstFile = files[0] || {};
    const url = firstFile.url || resource.downloadUrl;
    if (!url) {
      wx.showToast({ title: '暂无可打开的文件', icon: 'none' });
      return;
    }

    if (resource.fileType === '图片') {
      const imageUrls = files.filter((file) => getFileDisplayType(file) === '图片').map((file) => file.url);
      wx.previewImage({
        current: url,
        urls: imageUrls.length ? imageUrls : [url],
      });
      return;
    }

    wx.showLoading({ title: '打开中' });
    wx.downloadFile({
      url,
      success: (res) => {
        if (res.statusCode !== 200) {
          wx.showToast({ title: '下载失败', icon: 'none' });
          return;
        }
        wx.openDocument({
          filePath: res.tempFilePath,
          showMenu: true,
          fail: () => wx.showToast({ title: '暂不支持打开该文件', icon: 'none' }),
        });
      },
      fail: () => wx.showToast({ title: '下载失败', icon: 'none' }),
      complete: () => wx.hideLoading(),
    });
  },
});
