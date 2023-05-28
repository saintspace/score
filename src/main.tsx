import React from 'react'
import ReactDOM from 'react-dom/client'
import './index.css'
import {
  createBrowserRouter,
  RouterProvider,
} from "react-router-dom";
import { MantineProvider } from '@mantine/core';
import Root from "./pages/layouts/root";
import ErrorPage from "./pages/error/error";
import HomePage from "./pages/home/HomePage";
import AccountPage from "./pages/universe/account/AccountPage";
import DeleteAccountPage from './pages/universe/account/delete/DeleteAccountPage';
import SignOutPage from './pages/universe/auth/signout/SignOutPage';
import AuthCallbackPage from './pages/universe/auth/callback/AuthCallbackPage';
import VerifySubscriptionPage from './pages/universe/verify-subscription/VerifySubscriptionPage';

const router = createBrowserRouter([
  {
    path: "/",
    element: <Root />,
    errorElement: <ErrorPage />,
    children: [
      { index: true, element: <HomePage /> },
      { path: "saintspace/universe/account", element: <AccountPage /> },
      { path: "saintspace/universe/account/delete", element: <DeleteAccountPage /> },
      { path: "saintspace/universe/auth/signout", element: <SignOutPage /> },
      { path: "saintspace/universe/auth/callback", element: <AuthCallbackPage /> },
      { path: "saintspace/universe/verify-subscription", element: <VerifySubscriptionPage /> },
    ],
  },
]);

ReactDOM.createRoot(document.getElementById('root') as HTMLElement).render(
  <React.StrictMode>
    <MantineProvider withGlobalStyles withNormalizeCSS>
      <RouterProvider router={router} />
    </MantineProvider>
  </React.StrictMode>,
)
