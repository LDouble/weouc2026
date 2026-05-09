export const RESOURCE_CATEGORIES = [
  { label: '全部', value: 'all' },
  { label: '课程资料', value: 'course', icon: 'book' },
  { label: '考试经验', value: 'exam', icon: 'edit-1' },
  { label: '实验报告', value: 'lab', icon: 'file-copy' },
  { label: '办事指南', value: 'guide', icon: 'root-list' },
];

export const PUBLISH_CATEGORIES = RESOURCE_CATEGORIES.filter((item) => item.value !== 'all');
