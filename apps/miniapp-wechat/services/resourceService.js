import { fetchResourceDetail, fetchResourceList } from '~/api/modules/resource';
import { RESOURCE_CATEGORIES } from '~/constants/campus';
import { dedupeById, unwrapPayload } from './shared';

const CATEGORY_LABEL_MAP = RESOURCE_CATEGORIES.reduce((map, item) => {
  if (item.value && item.value !== 'all') {
    map[item.value] = item.label;
  }

  return map;
}, {});

const AVATAR_COLORS = ['blue', 'green', 'orange', 'purple', 'teal', 'rose'];

function stableIndex(value, modulo) {
  const text = `${value || ''}`;
  let total = 0;

  for (let index = 0; index < text.length; index += 1) {
    total += text.charCodeAt(index);
  }

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
    categoryLabel: CATEGORY_LABEL_MAP[category] || category || '其他',
    fileType,
    fileSize: firstFile.file_size || extra.file_size || '',
    files,
    downloadUrl: firstFile.url || extra.download_url || '',
    updatedAt: item.created_at || '',
    sponsor: item.publisher || '',
    avatarColor: AVATAR_COLORS[stableIndex(item.id, AVATAR_COLORS.length)],
    sponsorTag: extra.course_name || item.created_at || '',
    sponsorInitial: item.publisher_initial || (item.publisher ? item.publisher.charAt(0) : ''),
    views: Number(extra.views || 0),
    likes: Number(extra.likes || 0),
    fileIcon: getFileIcon(fileType),
  };
}

function moveItemToTop(items, id) {
  if (!id) return items;
  const target = items.find((item) => item.id === id);
  if (!target) return items;
  return [target].concat(items.filter((item) => item.id !== id));
}

export async function loadResourceList(params = {}) {
  const response = await fetchResourceList(params);
  const data = unwrapPayload(response);

  return {
    items: (data.list || []).map(mapResourceItem),
    total: Number(data.total || 0),
  };
}

export async function ensureResourcePinned(items, insertId) {
  if (!insertId) return items;

  const moved = moveItemToTop(items, insertId);
  if (moved[0] && moved[0].id === insertId) {
    return moved;
  }

  try {
    const response = await fetchResourceDetail(insertId);
    const detail = mapResourceItem(unwrapPayload(response));
    return [detail].concat(items.filter((item) => item.id !== insertId));
  } catch (error) {
    return items;
  }
}

export function mergeResourceItems(currentItems, nextItems) {
  return dedupeById([].concat(currentItems || [], nextItems || []));
}
