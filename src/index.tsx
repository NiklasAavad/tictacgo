import React from 'react';
import ReactDOM from 'react-dom/client';
import { Frontpage } from './components/Frontpage';
import './index.css';

const root = ReactDOM.createRoot(
  document.getElementById('root') as HTMLElement
);
root.render(
  <React.StrictMode>
    <Frontpage/>
  </React.StrictMode>
);