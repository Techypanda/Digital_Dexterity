import { useMemo } from "react";
import { useTokens } from "./utils";
import { decodeJwt } from 'jose'
export function useUsername(): string {
  const tokens = useTokens();
  return useMemo(() => {
    return decodeJwt(tokens.accessToken).username as string;
  }, [])
}