import { BrowserRouter, Routes, Route, Navigate } from 'react-router-dom';
import { AdminLayout } from '@/layout/Layout';
import { Login } from '@/pages/Login';
import { Dashboard } from '@/pages/Dashboard';
import { RoomManage } from '@/pages/RoomManage';
import { DeviceManage } from '@/pages/DeviceManage';
import { GroupManage } from '@/pages/GroupManage';
import { RuleManage } from '@/pages/RuleManage';
import { OrderManage } from '@/pages/OrderManage';
import { MerchManage } from '@/pages/MerchManage';
import { DeviceLogManage } from '@/pages/DeviceLogManage';

const ProtectedRoute = ({ children }: { children: React.ReactNode }) => {
  const token = localStorage.getItem('token');
  if (!token) {
    return <Navigate to="/login" />;
  }
  return <>{children}</>;
};

function App() {
  return (
    <BrowserRouter
      future={{
        v7_relativeSplatPath: true,
        v7_startTransition: true
      }}
    >
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <Dashboard />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/dashboard"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <Dashboard />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/room"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <RoomManage />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/device"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <DeviceManage />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/group"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <GroupManage />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/rule"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <RuleManage />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/order"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <OrderManage />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/merch"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <MerchManage />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
        <Route
          path="/devicelog"
          element={
            <ProtectedRoute>
              <AdminLayout>
                <DeviceLogManage />
              </AdminLayout>
            </ProtectedRoute>
          }
        />
      </Routes>
    </BrowserRouter>
  );
}

export default App;