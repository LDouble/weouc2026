export const CARPOOL_TIME_FILTERS = [
  { value: 'all', label: '全部' },
  { value: 'today', label: '今天' },
  { value: 'tomorrow', label: '明天' },
  { value: 'week', label: '本周' },
  { value: 'longterm', label: '长期' },
];

export const CARPOOL_CATEGORIES = [
  {
    value: 'today',
    label: '今日出发',
    icon: 'time',
    typeLabel: '今日顺路',
    emptyTitle: '今天还没有新的拼车',
    emptyDesc: '可以先发布行程，等顺路同学来拼单。',
  },
  {
    value: 'return',
    label: '返校',
    icon: 'home',
    typeLabel: '返校专线',
    emptyTitle: '暂无返校拼车',
    emptyDesc: '试试切换到出校或长期拼车，看看有没有合适路线。',
  },
  {
    value: 'leave',
    label: '出校',
    icon: 'location',
    typeLabel: '出校拼车',
    emptyTitle: '暂无出校行程',
    emptyDesc: '先留下你的出发时间和路线，系统会优先展示给顺路同学。',
  },
  {
    value: 'longterm',
    label: '长期拼车',
    icon: 'calendar',
    typeLabel: '长期通勤',
    emptyTitle: '长期路线还在招募中',
    emptyDesc: '发布常驻路线后，更容易匹配到固定搭子。',
  },
];
