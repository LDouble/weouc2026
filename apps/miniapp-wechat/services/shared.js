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
