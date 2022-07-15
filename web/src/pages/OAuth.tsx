import {Box, CircularProgress, Typography} from '@mui/material';
import {useContext, useEffect} from 'react';
import {useNavigate} from 'react-router-dom';
import {authContext, useGithubRedirect} from '../api/auth';

export default function OAuth() {
  const performAuth = useGithubRedirect();
  const navigate = useNavigate();
  const {setAuthenticated} = useContext(authContext);
  async function doGithubLogic() {
    const tokens = await performAuth();
    if (tokens) {
      localStorage.setItem('token', tokens.token);
      localStorage.setItem('refreshToken', tokens.refreshToken);
      setAuthenticated(true);
      navigate('/');
    } else {
      navigate('/');
    }
  }
  useEffect(() => {
    doGithubLogic();
  }, []);
  return (
    <Box display="flex" alignItems="center" justifyContent="center" width="100vw" height="100vh">
      <Box>
        <Box display="flex" justifyContent="center" mb={2}>
          <CircularProgress />
        </Box>
        <Typography>Performing Auth</Typography>
      </Box>
    </Box>
  );
}
