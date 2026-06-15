import React from 'react';
import { Shield, Activity, List, Settings } from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';

const Layout = ({ children }) => {
  const location = useLocation();

  const navItems = [
    { name: 'Dashboard', icon: Activity, path: '/' },
    { name: 'Audit Logs', icon: List, path: '/audit' },
    { name: 'Policies', icon: Shield, path: '/policies' },
    { name: 'Settings', icon: Settings, path: '/settings' },
  ];

  return (
    <div className="min-h-screen flex flex-col bg-surface">
      {/* Header */}
      <header className="bg-primary px-lg py-sm text-secondary flex items-center justify-between shadow-lg sticky top-0 z-50">
        <div className="flex items-center gap-3">
          <div className="bg-white p-2 rounded-md">
            <Shield className="text-primary w-6 h-6" />
          </div>
          <span className="text-xl font-bold tracking-tight">Aegis-MCP</span>
        </div>
        <nav className="hidden md:flex items-center gap-md">
          {navItems.map((item) => (
            <Link
              key={item.name}
              to={item.path}
              className={`flex items-center gap-2 px-3 py-2 rounded-md transition-colors ${
                location.pathname === item.path ? 'bg-white/20' : 'hover:bg-white/10'
              }`}
            >
              <item.icon className="w-4 h-4" />
              <span className="text-sm font-semibold">{item.name}</span>
            </Link>
          ))}
        </nav>
        <div className="flex items-center gap-4">
          <div className="w-8 h-8 rounded-full bg-primary-40 flex items-center justify-center font-bold text-xs uppercase">
            AD
          </div>
        </div>
      </header>

      {/* Main Content */}
      <main className="flex-1 container mx-auto px-lg py-xl">
        {children}
      </main>

      {/* Footer */}
      <footer className="bg-neutral border-t border-border/50 py-lg px-lg text-center text-sm text-tertiary">
        <p>&copy; 2026 Aegis-MCP Security Gateway. All rights reserved.</p>
      </footer>
    </div>
  );
};

export default Layout;
