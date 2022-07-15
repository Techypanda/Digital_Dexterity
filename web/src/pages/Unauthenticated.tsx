import {Alert, Box, Button, Divider, IconButton, Paper, TextField, Typography} from '@mui/material';
import {CommonProps} from '@mui/material/OverridableComponent';
import {useContext, useState} from 'react';
import {authContext, useGithubAuthRequest, useAuth} from '../api/auth';
import {HTTPError} from 'ky';
import {GitHub} from '@mui/icons-material';

export function Unauthenticated(props: CommonProps) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const {login, register} = useAuth();
  const {setAuthenticated} = useContext(authContext);
  const beginGithubAuth = useGithubAuthRequest();
  async function doLogin() {
    try {
      setLoading(true);
      setError('');
      setSuccess('');
      const loginResp = await (await login(username, password)).json<{ token: string; refreshToken: string; username: string }>();
      setSuccess(`welcome ${loginResp.username}`);
      localStorage.setItem('token', loginResp.token);
      localStorage.setItem('refreshToken', loginResp.refreshToken);
      setAuthenticated(true);
    } catch (err) {
      const resp = (await (err as HTTPError).response.json());
      setError(resp.error);
    } finally {
      setLoading(false);
    }
  }
  async function doRegister() {
    try {
      setLoading(true);
      setError('');
      setSuccess('');
      await register(username, password);
      setSuccess('you have successfully created a account');
    } catch (err) {
      const resp = (await (err as HTTPError).response.json());
      setError(resp.error);
    } finally {
      setLoading(false);
    }
  }
  return (
    <Box display="flex" width="100vw" height="100vh" justifyContent="center" alignItems="center">
      <Box maxWidth={600}>
        <Paper>
          <Typography variant="h4" component="h1" align="center" sx={{pt: 2}}>Digital Dexterity</Typography>
          <Typography variant="h6" component="h2" align="center" sx={{pb: 2}}>Sorry, This App Requires Usernames</Typography>
          <Divider />
          <Box p={2}>
            <TextField disabled={loading} variant="outlined" label="Username" value={username} onChange={(e) => setUsername(e.target.value)} fullWidth sx={{mb: 1}} />
            <TextField disabled={loading} variant="outlined" label="Password" value={password} type="password" autoComplete="current-password" onChange={(e) => setPassword(e.target.value)} fullWidth />
            {error && <Alert severity="error" sx={{mt: 1}}>{error}</Alert>}
            {success && <Alert severity="success" sx={{mt: 1}}>{success}</Alert>}
          </Box>
          <Divider />
          <Box p={2}>
            <Button disabled={loading} sx={{mr: 2}} onClick={() => doRegister()}>Register</Button>
            <Button disabled={loading} onClick={() => doLogin()}>Login</Button>
          </Box>
          <Divider />
          <Box p={2}>
            <Typography variant="subtitle2">Or Login With These Providers</Typography>
            <IconButton onClick={() => beginGithubAuth()} >
              <GitHub />
            </IconButton>
          </Box>
        </Paper>
      </Box>
    </Box>
  );
}
