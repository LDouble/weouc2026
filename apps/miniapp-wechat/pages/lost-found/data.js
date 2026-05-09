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

export const PUBLISH_CATEGORIES = LOST_FOUND_CATEGORIES.filter((item) => item.value !== 'all');
