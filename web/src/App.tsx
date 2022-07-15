import {ReactNode, useState} from 'react';
import {BrowserRouter, Route, Routes} from 'react-router-dom';
import {authContext, useAuthenticated} from './api/auth';
import {Explain} from './pages/Explain';
import Landing from './pages/Landing';
import Page from './pages/Page';
import Assess from './pages/Assess';
import {Unauthenticated} from './pages/Unauthenticated';
import AssessSelf from './pages/AssessSelf';
import OAuth from './pages/OAuth';

export function AuthenticationContext(props: { children: ReactNode }) {
  const [authenticated, setAuthenticated] = useState(false);
  return (
    <authContext.Provider value={{authenticated, setAuthenticated: (auth: boolean) => setAuthenticated(auth)}}>
      {props.children}
    </authContext.Provider>
  );
}

function App() {
  const {authenticated} = useAuthenticated();
  return (
    <>
      {authenticated ?
        <BrowserRouter>
          <Routes>
            <Route path="/assess/self" element={<Page to={<AssessSelf />} />} />
            <Route path="/assess" element={<Page to={<Assess />} />} />
            <Route path="/explain" element={<Page to={<Explain />} />} />
            <Route path="*" element={<Page to={<Landing />} />} />
          </Routes>
        </BrowserRouter> :
        <BrowserRouter>
          <Routes>
            <Route path="/oauth" element={<OAuth />} />
            <Route path="*" element={<Unauthenticated />} />
          </Routes>
        </BrowserRouter>
      }
    </>
  );
}

export default App;
