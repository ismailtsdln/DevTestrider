import { LayoutDashboard, ListTodo, Settings, Activity } from 'lucide-react';
import clsx from 'clsx';

interface SidebarProps {
  activeTab: string;
  setActiveTab: (tab: string) => void;
}

export function Sidebar({ activeTab, setActiveTab }: SidebarProps) {
  const items = [
    { id: 'dashboard', icon: LayoutDashboard, label: 'Dashboard' },
    { id: 'tests', icon: ListTodo, label: 'Test Suites' },
    { id: 'coverage', icon: Activity, label: 'Coverage' },
    { id: 'settings', icon: Settings, label: 'Settings' },
  ];

  return (
    <aside className="w-64 bg-slate-900 border-r border-slate-800 flex flex-col">
      <div className="p-6">
        <div className="h-8 w-8 bg-gradient-to-br from-indigo-500 to-cyan-500 rounded-lg mb-2" />
        <span className="font-bold text-lg tracking-tight">DevTestrider</span>
      </div>
      
      <nav className="flex-1 px-4 space-y-1">
        {items.map((item) => (
          <button
            key={item.id}
            onClick={() => setActiveTab(item.id)}
            className={clsx(
              "flex items-center space-x-3 w-full px-4 py-3 rounded-lg transition-all duration-200 text-sm font-medium",
              activeTab === item.id 
                ? "bg-indigo-500/10 text-indigo-400 border border-indigo-500/20 shadow-sm shadow-indigo-500/10" 
                : "text-slate-400 hover:text-slate-200 hover:bg-slate-800/50"
            )}
          >
            <item.icon size={20} className={activeTab === item.id ? "stroke-indigo-400" : "stroke-current"} />
            <span>{item.label}</span>
          </button>
        ))}
      </nav>
      
      <div className="p-4 border-t border-slate-800">
        <div className="bg-slate-800/50 rounded-lg p-4">
          <p className="text-xs font-mono text-slate-500">v0.1.0-beta</p>
          <div className="mt-2 text-xs text-slate-400 flex items-center gap-2">
            <div className="w-2 h-2 rounded-full bg-green-500/20 border border-green-500/50" />
            System Healthy
          </div>
        </div>
      </div>
    </aside>
  );
}
