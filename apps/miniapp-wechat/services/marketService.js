import { fetchMarketList, favoriteMarket } from '~/api/modules/market';
import { getMarketCategories } from '~/stores/config';
import { splitAlternateColumns, unwrapPayload } from './shared';

function buildCategoryLabelMap() {
  const categories = getMarketCategories();

  return categories.reduce((map, item) => {
    if (item.value && item.value !== 'all') {
      map[item.value] = item.label;
    }
    return map;
  }, {});
}

function mapMarketItem(item, categoryLabelMap) {
  const extra = item.extra || {};
  const isWanted = extra.category === 'wanted';
  const imageUrls = extra.images && extra.images.length ? extra.images.slice() : [];

  if (!imageUrls.length && item.image) {
    imageUrls.push(item.image);
  }

  const images = imageUrls.map((url, index) => ({
    key: `img-${item.id}-${index}`,
    url,
    variant: 'center',
  }));

  return {
    id: item.id,
    category: extra.category || 'life',
    categoryName: categoryLabelMap[extra.category] || extra.category || '',
    mode: isWanted ? 'wanted' : '',
    title: item.title || '',
    price: extra.price || '0',
    originalPrice: extra.original_price || '',
    condition: extra.condition || '',
    seller: item.publisher || '',
    sellerInitial: item.publisher_initial || '',
    avatarTone: 'green',
    likes: Number(extra.likes || item.likes || 0),
    liked: Boolean(extra.is_favorited || item.liked),
    image: item.image || (images[0] ? images[0].url : ''),
    images,
    imageHeight: images.length ? 342 : 354,
    cardBg: 'blue',
    detail: item.desc || '',
    contact: extra.contact || '',
  };
}

export function getMarketPageConfig() {
  return {
    categoryList: getMarketCategories(),
  };
}

export async function loadMarketList(params = {}) {
  const response = await fetchMarketList(params);
  const data = unwrapPayload(response);
  const categoryLabelMap = buildCategoryLabelMap();
  const items = (data.list || []).map((item) => mapMarketItem(item, categoryLabelMap));

  return {
    items,
    columns: splitAlternateColumns(items),
    total: Number(data.total || 0),
  };
}

export async function toggleMarketFavorite(product) {
  const action = product.liked ? 'remove' : 'add';
  await favoriteMarket(product.id, action);

  const liked = !product.liked;

  return {
    ...product,
    liked,
    likes: Math.max(0, Number(product.likes || 0) + (liked ? 1 : -1)),
  };
}
