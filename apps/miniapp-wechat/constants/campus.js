export const ERRAND_CATEGORIES = [
  { label: '全部', value: 'all' },
  { label: '取快递', value: 'parcel', icon: 'app' },
  { label: '代买饭', value: 'food', icon: 'shop' },
  { label: '帮占座', value: 'seat', icon: 'location' },
  { label: '代打印', value: 'print', icon: 'file-copy' },
  { label: '其他', value: 'other', icon: 'thunder' },
];

export const RESOURCE_CATEGORIES = [
  { label: '全部', value: 'all' },
  { label: '课程资料', value: 'course', icon: 'book' },
  { label: '考试经验', value: 'exam', icon: 'edit-1' },
  { label: '实验报告', value: 'lab', icon: 'file-copy' },
  { label: '办事指南', value: 'guide', icon: 'root-list' },
];

export const LOST_FOUND_TYPES = [
  { value: 'lost', label: '我丢了' },
  { value: 'found', label: '我捡到' },
];

export const LOST_FOUND_CATEGORIES = [
  { label: '全部', value: 'all' },
  { label: '证件卡片', value: 'card', icon: 'card' },
  { label: '电子设备', value: 'digital', icon: 'mobile' },
  { label: '书本文具', value: 'book', icon: 'books' },
  { label: '生活用品', value: 'life', icon: 'bag' },
  { label: '钥匙雨伞', value: 'key', icon: 'lock-on' },
];

export const ERRAND_PUBLISH_CATEGORIES = ERRAND_CATEGORIES.filter((item) => item.value !== 'all');
export const RESOURCE_PUBLISH_CATEGORIES = RESOURCE_CATEGORIES.filter((item) => item.value !== 'all');
export const LOST_FOUND_PUBLISH_CATEGORIES = LOST_FOUND_CATEGORIES.filter((item) => item.value !== 'all');
