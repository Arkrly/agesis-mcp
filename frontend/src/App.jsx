import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import Layout from './components/Layout';
import Dashboard from './pages/Dashboard';
import Policies from './pages/Policies';

// Placeholder components for other routes
const AuditLogs = () => <div className="p-10 text-center text-tertiary">Audit Logs View (Under Construction)</div>;
const Settings = () => <div className="p-10 text-center text-tertiary">Settings View (Under Construction)</div>;

function App() {
  return (
    <Router>
      <Layout>
        <Routes>
          <Route path="/" element={<Dashboard />} />
          <Route path="/audit" element={<AuditLogs />} />
          <Route path="/policies" element={<Policies />} />
          <Route path="/settings" element={<Settings />} />
        </Routes>
      </Layout>
    </Router>
  );
}

export default App;
