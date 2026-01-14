

export const API_BASE_URL: string = import.meta.env.VITE_API_URL;

if (!API_BASE_URL) {
  throw new Error("VITE_API_URL is not defined");
}

export async function apiFetch(path: string, options = {}): Promise<Response> {
  return fetch(`${API_BASE_URL}${path}`, {
    ...options,
  });
}
