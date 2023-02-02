export const API_URL = import.meta.env.VITE_API_URL;

export function fetchAll() {
  return fetchJsonApi(`${API_URL}/api/all`, "GET");
}

export function fetchStatus(serviceId: number) {
  return fetchJsonApi(`${API_URL}/api/status`, "GET");
}

export function fetchJsonApi(url, method) {
  return new Promise((res, rej) => {
    fetch(url, { method: method })
      .then((resp) => {
        resp
          .json()
          .then((r) => {
            res(r);
          })
          .catch(function (err) {
            rej(err);
          });
      })
      .catch(() => {
        rej(null);
      });
  });
}
