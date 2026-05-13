export function unwrapPayload(response) {
  if (response && typeof response === 'object' && Object.prototype.hasOwnProperty.call(response, 'data')) {
    return response.data || {};
  }

  return response || {};
}

export function splitAlternateColumns(items = []) {
  return items.reduce(
    (columns, item, index) => {
      const target = index % 2 === 0 ? 'left' : 'right';
      columns[target].push(item);
      return columns;
    },
    { left: [], right: [] },
  );
}

export function dedupeById(items = []) {
  const seen = {};

  return items.filter((item) => {
    if (!item || !item.id || seen[item.id]) return false;
    seen[item.id] = true;
    return true;
  });
}

export function getReviewStatus(item = {}) {
  const extra = item.extra || {};
  return item.review_status || item.status || extra.review_status || extra.status || '';
}

const STATUS_DISPLAY_MAP = {
  reviewing: { label: '审核中', tone: 'amber' },
  published: { label: '已发布', tone: 'green' },
  rejected: { label: '审核未通过', tone: 'red' },
  offline: { label: '已下架', tone: 'red' },
  cancelled: { label: '已取消', tone: 'red' },
  accepted: { label: '已接单', tone: 'blue' },
  open: { label: '报名中', tone: 'green' },
  full: { label: '人数已满', tone: 'purple' },
  resolved: { label: '已找到', tone: 'green' },
};

export function getStatusDisplay(status, overrides = {}) {
  const entry = STATUS_DISPLAY_MAP[status];
  if (!entry) return { label: status, tone: 'default' };
  return {
    label: (overrides[status] && overrides[status].label) || entry.label,
    tone: (overrides[status] && overrides[status].tone) || entry.tone,
  };
}

export function normalizeContactFields(item = {}) {
  const extra = item.extra || {};
  const rawCanViewContact = Object.prototype.hasOwnProperty.call(item, 'can_view_contact')
    ? item.can_view_contact
    : item.canViewContact;
  const canViewContact = typeof rawCanViewContact === 'boolean' ? rawCanViewContact : null;
  const contact = item.contact || extra.contact || '';
  const explicitReason =
    item.contact_hidden_reason ||
    item.contactHiddenReason ||
    extra.contact_hidden_reason ||
    extra.contactHiddenReason ||
    '';
  let contactHiddenReason = explicitReason;

  if (!contactHiddenReason && canViewContact === false) {
    contactHiddenReason = 'bind_required';
  } else if (!contactHiddenReason && !contact) {
    contactHiddenReason = 'empty';
  }

  return {
    canViewContact,
    contact,
    contactHiddenReason,
  };
}
