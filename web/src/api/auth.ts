import ky from "ky";
import { createContext, useContext, useEffect, useState } from "react";
import { useAPIURI } from "./utils";

const REFRESH_LOOP = 900000 // 15 Minutes

export async function refreshTokens() {
  if (localStorage.getItem("token") && localStorage.getItem("refreshToken")) {
    const resp = await ky(`${import.meta.env.VITE_API_URI}refresh`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem("refreshToken")}`,
      }
    }).json<{ token: string; refreshToken: string }>();
    localStorage.setItem("token", resp.token);
    localStorage.setItem("refreshToken", resp.refreshToken);
  }
}

export function useAuthenticated() {
  const { authenticated, setAuthenticated } = useContext(authContext);
  const [refreshLoop, setRefreshLoop] = useState<NodeJS.Timer>();
  async function validateTokens() {
    try {
      await refreshTokens()
      if (refreshLoop) {
        clearInterval(refreshLoop);
      }
      setRefreshLoop(setInterval(() => {
        refreshTokens();
      }, REFRESH_LOOP))
      setAuthenticated(true);
    } catch (err) {
    }
  }
  useEffect(() => {
    const token = localStorage.getItem("token");
    const refreshToken = localStorage.getItem("refreshToken");
    if (token && refreshToken) {
      validateTokens();
    }
  }, [])
  return { authenticated, validateTokens };
}

export function logout() {
  localStorage.clear();
}

export function useAuth() {
  const apiURI = useAPIURI();
  async function login(username: string, password: string) {
    return ky(`${apiURI}login`, {
      json: { username, password },
      method: 'post',
    })
  }
  async function register(username: string, password: string) {
    return ky(`${apiURI}register`, {
      json: { username, password },
      method: 'post',
    })
  }
  return { login, register };
}

export const authContext = createContext({ authenticated: false, setAuthenticated: (auth: boolean) => {} })