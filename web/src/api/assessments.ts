import ky from "ky";
import { useMutation, useQuery, useQueryClient } from "react-query";
import { useAPIURI, useTokens } from "./utils";

export function useAssessments() {
  const baseAPIURI = useAPIURI();
  const tokens = useTokens();
  return useQuery(['assessments'], async () => {
    return await ky(`${baseAPIURI}api/v1/assess`, {
      headers: {
        Authorization: `Bearer ${tokens.accessToken}`
      },
      method: 'get',
    }).json<AssessmentsResponse>()
  }, {
    cacheTime: 30000,
    staleTime: 30000,
  })
}

export function useNewAssessment() {
  const baseAPIURI = useAPIURI();
  const tokens = useTokens();
  const client = useQueryClient();
  const selfAssessment = useMutation(async (assessment: Assessment) => {
    return await ky(`${baseAPIURI}api/v1/assess/me`, {
      headers: {
        Authorization: `Bearer ${tokens.accessToken}`
      },
      method: 'post',
      json: assessment,
    }).json();
  }, {
    onSettled: () => client.invalidateQueries('assessments')
  });
  const externalAssessment = useMutation(async (assessment: Assessment & { assessing: string }) => {
    return await ky(`${baseAPIURI}api/v1/assess`, {
      headers: {
        Authorization: `Bearer ${tokens.accessToken}`
      },
      method: 'post',
      json: assessment,
    }).json();
  }, {
  });
  return { selfAssessment, externalAssessment }
}