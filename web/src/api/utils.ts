import { useMemo } from "react";

export function useAPIURI() {
  return useMemo(() => import.meta.env.VITE_API_URI, []);
}
export function useTokens() {
  return { accessToken: localStorage.getItem("token")!, refreshToken: localStorage.getItem("refreshToken")! };
}