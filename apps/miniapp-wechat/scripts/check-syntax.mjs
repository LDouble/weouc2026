import { execFileSync } from 'node:child_process';
import { readdirSync, statSync } from 'node:fs';
import { resolve, extname } from 'node:path';

const ROOT = resolve(process.cwd());
const TARGETS = [
  'app.js',
  'api',
  'services',
  'stores',
  'pages',
  'components',
  'behaviors',
  'utils',
];
const EXCLUDE_DIRS = new Set(['node_modules', 'miniprogram_npm']);

function collectJSFiles(pathname, files) {
  const stats = statSync(pathname);
  if (stats.isFile()) {
    if (extname(pathname) === '.js') files.push(pathname);
    return;
  }

  if (!stats.isDirectory()) return;

  const entries = readdirSync(pathname, { withFileTypes: true });
  for (const entry of entries) {
    if (entry.isDirectory() && EXCLUDE_DIRS.has(entry.name)) continue;
    collectJSFiles(resolve(pathname, entry.name), files);
  }
}

const jsFiles = [];
for (const target of TARGETS) {
  collectJSFiles(resolve(ROOT, target), jsFiles);
}
jsFiles.sort();

if (jsFiles.length === 0) {
  console.error('未找到可校验的 JS 文件');
  process.exit(1);
}

const failed = [];
for (const file of jsFiles) {
  try {
    execFileSync(process.execPath, ['--check', file], { stdio: 'pipe' });
  } catch (error) {
    failed.push({
      file,
      output: String(error.stderr || error.stdout || error.message || ''),
    });
  }
}

if (failed.length > 0) {
  console.error(`语法校验失败：${failed.length}/${jsFiles.length}`);
  for (const item of failed) {
    console.error(`\n[${item.file}]`);
    console.error(item.output.trim());
  }
  process.exit(1);
}

console.log(`语法校验通过：${jsFiles.length} 个文件`);
