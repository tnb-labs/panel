import { request } from '@/utils'

export default {
  // 设为 CLI 版本
  setCli: (version: number): any => request.post(`/apps/php${version}/setCli`),
  // 获取配置
  config: (version: number): any => request.get(`/apps/php${version}/config`),
  // 保存配置
  saveConfig: (version: number, config: string): any =>
    request.post(`/apps/php${version}/config`, { config }),
  // 获取FPM配置
  fpmConfig: (version: number): any => request.get(`/apps/php${version}/fpmConfig`),
  // 保存FPM配置
  saveFPMConfig: (version: number, config: string): any =>
    request.post(`/apps/php${version}/fpmConfig`, { config }),
  // 负载状态
  load: (version: number): any => request.get(`/apps/php${version}/load`),
  // 获取错误日志
  errorLog: (version: number): any => request.get(`/apps/php${version}/errorLog`),
  // 清空错误日志
  clearErrorLog: (version: number): any => request.post(`/apps/php${version}/clearErrorLog`),
  // 获取慢日志
  slowLog: (version: number): any => request.get(`/apps/php${version}/slowLog`),
  // 清空慢日志
  clearSlowLog: (version: number): any => request.post(`/apps/php${version}/clearSlowLog`),
  // 拓展列表
  extensions: (version: number): any => request.get(`/apps/php${version}/extensions`),
  // 安装拓展
  installExtension: (version: number, slug: string): any =>
    request.post(`/apps/php${version}/extensions`, { slug }),
  // 卸载拓展
  uninstallExtension: (version: number, slug: string): any =>
    request.delete(`/apps/php${version}/extensions`, { data: { slug } })
}
