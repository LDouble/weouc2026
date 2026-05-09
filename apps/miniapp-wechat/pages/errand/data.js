export const ERRAND_CATEGORIES = [
  { label: '全部', value: 'all' },
  { label: '取快递', value: 'parcel', icon: 'app' },
  { label: '代买饭', value: 'food', icon: 'shop' },
  { label: '帮占座', value: 'seat', icon: 'location' },
  { label: '代打印', value: 'print', icon: 'file-copy' },
  { label: '其他', value: 'other', icon: 'thunder' },
];

export const PUBLISH_CATEGORIES = ERRAND_CATEGORIES.filter((item) => item.value !== 'all');
