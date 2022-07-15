import React from 'react';
import ReactDOM from 'react-dom/client';
import {QueryClientProvider, QueryClient} from 'react-query';
import {ReactQueryDevtools} from 'react-query/devtools';
import App, {AuthenticationContext} from './App';
import './assets/style.css';

const client = new QueryClient();
ReactDOM.createRoot(document.getElementById('root')!).render(
    <React.StrictMode>
      <AuthenticationContext>
        <QueryClientProvider client={client}>
          <App />
          <ReactQueryDevtools initialIsOpen={false} />
        </QueryClientProvider>
      </AuthenticationContext>
    </React.StrictMode>,
);
