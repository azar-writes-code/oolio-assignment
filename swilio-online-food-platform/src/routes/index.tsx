import { BasicLayout } from "@/components/layouts";
import { Suspense } from "react";
import { Navigate, Outlet, useRoutes } from "react-router-dom";
import * as Pages from "./paths";

export default function AppRouter() {
  const { HomePage } = Pages;

  const publicRoutes = [
    {
      path: "/",
      element: (
        <BasicLayout>
          <Suspense fallback={<div>Loading...</div>}>
            <Outlet />
          </Suspense>
        </BasicLayout>
      ),
      children: [
        {
          index: true,
          element: <HomePage />,
        },
      ],
    },
    {
      path: "*",
      element: <Navigate to="/404" replace />,
    },
  ];

  const routes = useRoutes([...publicRoutes]);

  return routes;
}
