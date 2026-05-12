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
