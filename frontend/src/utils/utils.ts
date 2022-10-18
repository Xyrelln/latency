/**
 * 检查是否运行在 wails 环境
 */
export function isWailsRun(){
  // @ts-ignore: wails 运行时检查
  return window.runtime !== undefined
}