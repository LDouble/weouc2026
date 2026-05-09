# 公共业务页组件规范

本目录用于存放跨业务页面复用的基础 UI 组件。适用范围包括跑腿、交易、资料、拼车、失物招领等业务列表页中的顶部栏、搜索框、筛选、空状态和悬浮发布入口。

## 目录与引用

- 公共组件统一放在 `components/common/<component-name>/`。
- 组件文件统一使用 `index.json`、`index.wxml`、`index.less`、`index.js`。
- 页面必须在 `index.json` 中使用根路径引用：

```json
{
  "usingComponents": {
    "common-page-header": "/components/common/common-page-header/index",
    "common-search-bar": "/components/common/common-search-bar/index",
    "common-filter-tabs": "/components/common/common-filter-tabs/index",
    "common-empty-state": "/components/common/common-empty-state/index",
    "common-fab": "/components/common/common-fab/index"
  }
}
```

## 组件边界

- `components/common/` 只放跨业务 UI，不放商品、跑腿任务、资料、拼车行程、失物招领卡片等强业务组件。
- 页面继续持有业务状态，例如 `searchKeyword`、`activeCategory`、`visibleItems`。
- 公共组件只负责展示、基础交互和统一事件派发，不直接执行业务筛选、跳转或数据变更。
- 业务主题色、hero 区、列表布局和卡片内容保留在页面层。

## 组件约定

### common-page-header

用于 `navigationStyle: "custom"` 页面顶部区域。

常用属性：

- `title`：标题。
- `subtitle`：副标题，可选。
- `right-icon`：右侧图标，可选。
- `background`、`border-color`、`button-background`：视觉定制。

事件：

- `bind:back`：点击返回入口。
- `bind:righttap`：点击右侧操作入口。
- `bind:heightchange`：顶部栏高度变化，事件参数为 `{ height }`，页面应使用该值设置内容区 `padding-top`。

### common-search-bar

受控搜索框组件。

常用属性：

- `value`：当前搜索值。
- `placeholder`：占位文案。
- `clearable`：是否显示清空入口。
- `variant="raised"`：用于需要浮起阴影的搜索框。

事件：

- `bind:input`：输入变化，事件参数为 `{ value }`。
- `bind:confirm`：确认搜索，事件参数为 `{ value }`。
- `bind:clear`：清空搜索，事件参数为 `{ value: "" }`。
- `bind:tap`：点击搜索框容器。

### common-filter-tabs

受控筛选组件。

常用属性：

- `items`：筛选项数组。
- `active`：当前选中值。
- `variant`：`pill`、`line`、`segment`。
- `scrollable`：是否横向滚动。
- `value-key`、`label-key`：自定义取值字段和展示字段。
- `active-color`：选中态主色。

事件：

- `bind:change`：切换筛选项，事件参数为 `{ value, item, index }`。
- 点击当前 active 项不会重复触发 `change`。

### common-empty-state

用于空列表、无搜索结果、暂无数据等状态。

常用属性：

- `icon`：TDesign 图标名。
- `title`：标题。
- `description`：说明文案。
- `action-text`：操作按钮文案，可选。

事件：

- `bind:action`：点击操作按钮，由页面执行发布、上传、登记或重试等业务动作。

### common-fab

用于页面悬浮发布入口。

常用属性：

- `icon`：默认 `add`。
- `text`：可选文案；有文案时为扩展按钮，无文案时为圆形按钮。
- `background`、`shadow`：视觉定制。

事件：

- `bind:tap`：点击悬浮按钮，由页面执行对应业务跳转。

## 新增业务页接入要求

新增与跑腿、交易、资料、拼车、失物招领相似的业务列表页时，应优先使用：

- `common-page-header` 处理自定义顶部栏和安全区。
- `common-search-bar` 处理搜索输入、确认和清空。
- `common-filter-tabs` 处理分类、状态、时间等筛选。
- `common-empty-state` 处理空列表和无搜索结果。
- `common-fab` 处理发布、上传、登记等悬浮入口。

只有当页面存在强业务差异，且无法通过属性、插槽或页面层样式表达时，才新增业务专属组件。
