import createStore from './createStore';

const MARKET_CATEGORIES = [
  { label: '全部闲置', value: 'all' },
  { label: '数码电子', value: 'digital' },
  { label: '教材教辅', value: 'book' },
  { label: '交通代步', value: 'transport' },
  { label: '生活用品', value: 'life' },
  { label: '服饰美妆', value: 'clothing' },
  { label: '求购代办', value: 'wanted' },
];

const RESOURCE_CATEGORIES = [
  { name: '课程资料', icon: 'book', count: 128 },
  { name: '考试经验', icon: 'edit-1', count: 46 },
  { name: '实验报告', icon: 'file-copy', count: 73 },
  { name: '办事指南', icon: 'root-list', count: 21 },
];

function getCategoryLabels(categories) {
  return categories.filter((item) => item.value !== 'all' && item.label).map((item) => item.label);
}

function getResourceCategoryNames(categories) {
  return categories.filter((item) => item.name).map((item) => item.name);
}

function createCategoryValueMap(categories) {
  return categories.reduce((map, item) => {
    if (item.value && item.value !== 'all' && item.label) {
      map[item.label] = item.value;
    }

    return map;
  }, {});
}

const RELEASE_SCENES = [
  {
    key: 'market',
    title: '发布闲置',
    descPlaceholder: '描述一下你要出手的闲置物品，比如品牌、型号、成色、转手原因...',
    priceLabel: '出手价格',
    categoryTitle: '分类',
    conditionTitle: '物品成色',
    tradeTitle: '交易方式',
    tradeMode: '支持校园面交',
    contactPlaceholder: '联系方式（微信号、手机号或站内私信）',
    showOriginal: true,
    showCardBg: true,
    categories: getCategoryLabels(MARKET_CATEGORIES),
    conditions: ['全新', '99新', '9成新', '8成新', '战损版', '求购'],
    extraFields: [],
    defaults: {
      category: '数码电子',
      condition: '99新',
    },
  },
  {
    key: 'errand',
    title: '发布跑腿',
    descPlaceholder: '说明需要同学帮忙做什么、取送地点、截止时间和注意事项...',
    priceLabel: '任务酬劳',
    categoryTitle: '任务类型',
    conditionTitle: '紧急程度',
    tradeTitle: '服务范围',
    tradeMode: '校内可接单',
    contactPlaceholder: '联系方式（便于接单同学联系你）',
    showOriginal: false,
    showCardBg: false,
    categories: ['代取快递', '代送物品', '排队占位', '打印帮取', '其他帮办'],
    conditions: ['不急', '今天内', '1小时内', '越快越好'],
    extraFields: [
      { field: 'startPlace', label: '起点', icon: 'location', placeholder: '东门快递驿站' },
      { field: 'endPlace', label: '终点', icon: 'send', placeholder: '7号楼楼下' },
      { field: 'deadline', label: '截止', icon: 'time', placeholder: '今天 19:00 前' },
    ],
    defaults: {
      category: '代取快递',
      condition: '今天内',
    },
  },
  {
    key: 'lostFound',
    title: '发布登记',
    descPlaceholder: '描述物品特征、丢失或拾取地点，避免写出可被冒领的完整信息...',
    priceLabel: '悬赏金额',
    categoryTitle: '登记类型',
    conditionTitle: '物品类别',
    tradeTitle: '联系地点',
    tradeMode: '校内联系',
    contactPlaceholder: '联系方式（便于认领或找回）',
    showOriginal: false,
    showCardBg: true,
    categories: ['我丢了', '我捡到'],
    conditions: ['证件卡片', '电子设备', '书本文具', '生活用品', '钥匙雨伞'],
    extraFields: [
      { field: 'eventTime', label: '时间', icon: 'time', placeholder: '今天 12:30 左右' },
      { field: 'itemFeature', label: '特征', icon: 'tag', placeholder: '蓝色卡套、背面有贴纸' },
    ],
    defaults: {
      category: '我丢了',
      condition: '证件卡片',
    },
  },
  {
    key: 'resource',
    title: '上传资料',
    descPlaceholder: '说明资料内容、适用课程、年份范围和获取方式...',
    priceLabel: '资料价格',
    categoryTitle: '资料类型',
    conditionTitle: '文件格式',
    tradeTitle: '获取方式',
    tradeMode: '站内获取',
    contactPlaceholder: '联系方式或资料获取说明',
    showOriginal: false,
    showCardBg: true,
    categories: getResourceCategoryNames(RESOURCE_CATEGORIES),
    conditions: ['PDF', 'Word', '图片截图', '网盘链接'],
    extraFields: [
      { field: 'courseName', label: '课程', icon: 'file-copy', placeholder: '高等数学 A / 转专业申请' },
      { field: 'resourceScope', label: '范围', icon: 'tag', placeholder: '2025 秋季、计算机学院' },
    ],
    defaults: {
      category: '课程资料',
      condition: 'PDF',
    },
  },
];

const PALETTE_OPTIONS = [
  { key: 'purple', background: 'linear-gradient(135deg, #615fff, #9810fa)' },
  { key: 'rose', background: 'linear-gradient(135deg, #ff637e, #fb2c36)' },
  { key: 'orange', background: 'linear-gradient(135deg, #ffb900, #ff6900)' },
  { key: 'green', background: 'linear-gradient(135deg, #00d5be, #00bc7d)' },
  { key: 'blue', background: 'linear-gradient(135deg, #51a2ff, #00b8db)' },
  { key: 'dark', background: 'linear-gradient(135deg, #314158, #0f172b)' },
];

export const DEFAULT_APP_CONFIG = {
  market: {
    categories: MARKET_CATEGORIES,
    categoryValueMap: createCategoryValueMap(MARKET_CATEGORIES),
  },
  release: {
    scenes: RELEASE_SCENES,
    paletteOptions: PALETTE_OPTIONS,
  },
  resource: {
    categories: RESOURCE_CATEGORIES,
  },
};

const appConfigStore = createStore({
  appConfig: DEFAULT_APP_CONFIG,
  loading: false,
  loaded: false,
  loadedAt: 0,
  error: '',
});

function hasList(value) {
  return Array.isArray(value) && value.length > 0;
}

function normalizeConfig(config = {}) {
  const source = config || {};
  const market = source.market || {};
  const release = source.release || {};
  const resource = source.resource || {};
  const marketCategories = hasList(market.categories) ? market.categories : DEFAULT_APP_CONFIG.market.categories;
  const resourceCategories = hasList(resource.categories) ? resource.categories : DEFAULT_APP_CONFIG.resource.categories;
  const releaseScenes = hasList(release.scenes)
    ? release.scenes
    : syncReleaseSceneCategories(DEFAULT_APP_CONFIG.release.scenes, marketCategories, resourceCategories);

  return {
    market: {
      categories: marketCategories,
      categoryValueMap: Object.assign(
        {},
        DEFAULT_APP_CONFIG.market.categoryValueMap,
        createCategoryValueMap(marketCategories),
        market.categoryValueMap || {},
      ),
    },
    release: {
      scenes: releaseScenes,
      paletteOptions: hasList(release.paletteOptions) ? release.paletteOptions : DEFAULT_APP_CONFIG.release.paletteOptions,
    },
    resource: {
      categories: resourceCategories,
    },
  };
}

function syncReleaseSceneCategories(scenes, marketCategories, resourceCategories) {
  return scenes.map((scene) => {
    if (scene.key === 'market') {
      return Object.assign({}, scene, {
        categories: getCategoryLabels(marketCategories),
      });
    }

    if (scene.key === 'resource') {
      return Object.assign({}, scene, {
        categories: getResourceCategoryNames(resourceCategories),
      });
    }

    return scene;
  });
}

export function initAppConfigStore(initialConfig = DEFAULT_APP_CONFIG) {
  return setAppConfig(initialConfig);
}

export function getAppConfig() {
  return appConfigStore.getState().appConfig;
}

export function getAppConfigState() {
  return appConfigStore.getState();
}

export function setAppConfig(nextConfig) {
  const appConfig = normalizeConfig(nextConfig);
  return appConfigStore.setState({
    appConfig,
    loaded: true,
    loadedAt: Date.now(),
    error: '',
  }).appConfig;
}

export function subscribeAppConfig(listener) {
  if (typeof listener !== 'function') return () => {};

  return appConfigStore.subscribe((state, previousState) => {
    listener(state.appConfig, previousState.appConfig);
  });
}

export function loadAppConfig(loader) {
  if (typeof loader !== 'function') return Promise.resolve(getAppConfig());

  appConfigStore.setState({ loading: true, error: '' });

  return Promise.resolve()
    .then(loader)
    .then((remoteConfig) => {
      const appConfig = setAppConfig(remoteConfig || DEFAULT_APP_CONFIG);
      appConfigStore.setState({ loading: false });
      return appConfig;
    })
    .catch((error) => {
      appConfigStore.setState({
        loading: false,
        error: error && error.message ? error.message : '配置加载失败',
      });
      return getAppConfig();
    });
}

export function getMarketCategories() {
  return getAppConfig().market.categories;
}

export function getMarketCategoryValueMap() {
  return Object.assign({}, getAppConfig().market.categoryValueMap);
}

export function getReleaseScenes() {
  return getAppConfig().release.scenes;
}

export function getReleaseScene(type = 'market') {
  const scenes = getReleaseScenes();
  return scenes.find((item) => item.key === type) || scenes[0];
}

export function getPaletteOptions() {
  return getAppConfig().release.paletteOptions;
}

export function getResourceCategories() {
  return getAppConfig().resource.categories;
}
