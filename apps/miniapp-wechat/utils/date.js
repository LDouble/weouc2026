function padNumber(value) {
  return value < 10 ? `0${value}` : `${value}`;
}

export function safeParseDate(value) {
  if (!value && value !== 0) return null;
  if (value instanceof Date) {
    return Number.isNaN(value.getTime()) ? null : new Date(value.getTime());
  }

  const normalizedValue = typeof value === 'string' ? value.replace(/-/g, '/') : value;
  const date = new Date(normalizedValue);
  return Number.isNaN(date.getTime()) ? null : date;
}

export function formatDateTime(value) {
  const date = safeParseDate(value);
  if (!date) return '';

  return `${date.getFullYear()}-${padNumber(date.getMonth() + 1)}-${padNumber(date.getDate())} ${padNumber(date.getHours())}:${padNumber(date.getMinutes())}`;
}

export function formatRelativeTime(value) {
  const date = safeParseDate(value);
  if (!date) return '';

  const now = Date.now();
  const diff = now - date.getTime();

  if (diff < 60000) return '刚刚';
  if (diff < 3600000) return `${Math.floor(diff / 60000)} 分钟前`;
  if (diff < 86400000) return `${Math.floor(diff / 3600000)} 小时前`;
  if (diff < 604800000) return `${Math.floor(diff / 86400000)} 天前`;

  return `${date.getFullYear()}-${padNumber(date.getMonth() + 1)}-${padNumber(date.getDate())}`;
}

export function formatDeadline(value) {
  const date = safeParseDate(value);
  if (!date) return '';

  const now = new Date();
  const todayKey = `${now.getFullYear()}-${now.getMonth()}-${now.getDate()}`;
  const tomorrow = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1);
  const dateKey = `${date.getFullYear()}-${date.getMonth()}-${date.getDate()}`;
  const tomorrowKey = `${tomorrow.getFullYear()}-${tomorrow.getMonth()}-${tomorrow.getDate()}`;
  const clockText = `${padNumber(date.getHours())}:${padNumber(date.getMinutes())}`;

  if (dateKey === todayKey) return `今天 ${clockText} 前`;
  if (dateKey === tomorrowKey) return `明天 ${clockText} 前`;

  return `${date.getMonth() + 1}月${date.getDate()}日 ${clockText} 前`;
}
