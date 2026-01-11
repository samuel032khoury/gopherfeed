import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import { ActivationPage } from './ActivationPage.tsx'
import App from './App.tsx'
import './index.css'


const router = createBrowserRouter([
  {path: "/", element: <App />},
  {path: "/activate", element: <ActivationPage />},
    ])

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>,
)
