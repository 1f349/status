import { firstAvailableValue } from "./utils";

export function popupCenterScreen(url, title, w, h, focus) {
  const top = (screen.availHeight - h) / 4;
  const left = (screen.availWidth - w) / 2;
  const popup = openWindow(url, title, `scrollbars=yes,width=${w},height=${h},top=${top},left=${left}`);
  if (focus === true && window.focus) popup.focus();
  return popup;
}

export function openWindow(url, winnm, options) {
  const wTop = firstAvailableValue([window.screen.availTop, window.screenY, window.screenTop, 0]);
  const wLeft = firstAvailableValue([window.screen.availLeft, window.screenX, window.screenLeft, 0]);
  let top = 0;
  let left = 0;
  let result;
  if ((result = /top=(\d+)/g.exec(options))) top = parseInt(result[1]);
  if ((result = /left=(\d+)/g.exec(options))) left = parseInt(result[1]);
  let w;
  if (options) {
    options = options.replace("top=" + top, "top=" + (parseInt(top) + wTop));
    options = options.replace("left=" + left, "left=" + (parseInt(left) + wLeft));
    w = window.open(url, winnm, options);
  } else w = window.open(url, winnm);
  return w;
}
