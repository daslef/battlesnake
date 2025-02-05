import { createRoot } from 'react-dom/client'
import { Theme } from "@radix-ui/themes";
import App from './App.jsx'

import "@radix-ui/themes/styles.css";
import './index.css'

const rootElement = document.getElementById('root')

if (rootElement) {
  createRoot(rootElement).render(
    <Theme appearance="dark">
      <App />
    </Theme>
  )
}
