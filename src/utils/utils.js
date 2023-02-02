export function firstAvailableValue(arr) {
  for (let i = 0; i < arr.length; i++) if (typeof arr[i] != "undefined") return arr[i];
}

export function isObject(obj) {
  return obj != null && obj.constructor.name === "Object";
}
