import React, { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from "react-router";
import './index.css'
import '@fontsource/roboto/300.css';
import '@fontsource/roboto/400.css';
import '@fontsource/roboto/500.css';
import '@fontsource/roboto/700.css';
import {createTheme, CssBaseline, ThemeProvider} from "@mui/material";
import {MainWindow} from "./mainwindow.tsx";


const router = createBrowserRouter([
    {
        path: "/",
        element: <MainWindow/>,
    }
]);

const theme = createTheme();

createRoot(document.getElementById('root')!).render(
  <StrictMode>
      <CssBaseline>
          <ThemeProvider theme={theme}>
            <RouterProvider router={router}/>
          </ThemeProvider>
      </CssBaseline>
  </StrictMode>
)
