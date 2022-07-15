import ky from 'ky';
import {createContext, useContext, useEffect, useState} from 'react';
import {useNavigate, useSearchParams} from 'react-router-dom';
import {useAPIURI} from './utils';

const REFRESH_LOOP = 900000; // 15 Minutes

export async function refreshTokens() {
  if (localStorage.getItem('token') && localStorage.getItem('refreshToken')) {
    const resp = await ky(`${import.meta.env.VITE_API_URI}refresh`, {
      headers: {
        Authorization: `Bearer ${localStorage.getItem('refreshToken')}`,
      },
    }).json<{ token: string; refreshToken: string }>();
    localStorage.setItem('token', resp.token);
    localStorage.setItem('refreshToken', resp.refreshToken);
  }
}

export function useAuthenticated() {
  const {authenticated, setAuthenticated} = useContext(authContext);
  const [refreshLoop, setRefreshLoop] = useState<NodeJS.Timer>();
  async function validateTokens() {
    try {
      await refreshTokens();
      if (refreshLoop) {
        clearInterval(refreshLoop);
      }
      setRefreshLoop(setInterval(() => {
        refreshTokens();
      }, REFRESH_LOOP));
      setAuthenticated(true);
    } catch (err) {
    }
  }
  useEffect(() => {
    const token = localStorage.getItem('token');
    const refreshToken = localStorage.getItem('refreshToken');
    if (token && refreshToken) {
      validateTokens();
    }
  }, []);
  return {authenticated, validateTokens};
}

export function logout() {
  localStorage.clear();
}

export function useGithubAuthRequest() {
  const apiURI = useAPIURI();
  return async function beginGithubAuth() {
    const somewhatSecretState = `${Math.random().toString(36)}${Math.random().toString(36)}${Math.random().toString(36)}${Math.random().toString(36)}`;
    const githubURL = new URL('https://github.com/login/oauth/authorize');
    githubURL.searchParams.append('client_id', import.meta.env.VITE_GITHUB_CLIENT_ID);
    githubURL.searchParams.append('redirect_uri', import.meta.env.VITE_GITHUB_REDIRECT);
    githubURL.searchParams.append('scope', 'user');
    githubURL.searchParams.append('state', somewhatSecretState);
    await ky(`${apiURI}github/state`, {
      method: 'post',
      json: {
        state: somewhatSecretState,
      },
    });
    window.location.href = githubURL.toString();
  };
}

export function useGithubRedirect() {
  const [params] = useSearchParams();
  const navigate = useNavigate();
  const apiURI = useAPIURI();
  return async function doGithubAuthentication(): Promise<{ token: string, refreshToken: string } | null> {
    if (!params.get('code') || !params.get('state')) {
      navigate('/');
    }
    try {
      const resp = await ky(`${apiURI}github/login`, {
        method: 'post',
        json: {
          state: params.get('state')!,
          code: params.get('code')!,
          redirectURI: import.meta.env.VITE_GITHUB_REDIRECT,
        },
        timeout: 60000,
      }).json();
      return resp as { token: string, refreshToken: string };
    } catch (err) {
      console.error(err);
      navigate('/');
      return null;
    }
  };
}

export function useAuth() {
  const apiURI = useAPIURI();
  async function login(username: string, password: string) {
    return ky(`${apiURI}login`, {
      json: {username, password},
      method: 'post',
    });
  }
  async function register(username: string, password: string) {
    return ky(`${apiURI}register`, {
      json: {username, password},
      method: 'post',
    });
  }
  return {login, register};
}
// https:// github.com/login/oauth/authorize?client_id=7fdd431d8788fe30129b&redirect_uri=https://digitaldexterity.techytechster.com/oauth&scope=user&state=...
export const authContext = createContext({authenticated: false, setAuthenticated: (auth: boolean) => { }});
