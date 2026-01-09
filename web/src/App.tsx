import { useState } from 'react';
import { Sidebar } from './components/Sidebar';
import { Dashboard } from './components/Dashboard';
import { TestDetails } from './components/TestDetails';
import { Toaster } from 'react-hot-toast'; // We might need to install this or use a simple one

// Using a simple state manager or context would be good, but prop drilling is fine for this size
function App() {
  const [activeTab, setActiveTab] = useState('dashboard');

  return (
    <div className="flex h-screen bg-slate-950 text-slate-50 font-sans selection:bg-indigo-500/30">
      <Sidebar activeTab={activeTab} setActiveTab={setActiveTab} />
      
      <main className="flex-1 overflow-auto bg-slate-900/50 relative">
        <div className="absolute inset-0 bg-grid-slate-900/[0.04] bg-[bottom_1px_center] [mask-image:linear-gradient(to_bottom,transparent,black)] pointer-events-none" />
        
        <div className="relative p-8 max-w-7xl mx-auto">
          <header className="mb-8 flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold tracking-tight text-transparent bg-clip-text bg-gradient-to-r from-indigo-400 to-cyan-400">
                DevTestrider
              </h1>
              <p className="text-slate-400 mt-1">Real-time Test Intelligence</p>
            </div>
            <div className="flex items-center space-x-4">
              <span className="flex h-3 w-3 relative">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-emerald-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-emerald-500"></span>
              </span>
              <span className="text-sm font-medium text-emerald-400">Engine Active</span>
            </div>
          </header>

          {activeTab === 'dashboard' && <Dashboard />}
          {activeTab === 'tests' && <TestDetails />}
        </div>
      </main>
      <Toaster position="bottom-right" toastOptions={{ style: { background: '#1e293b', color: '#fff' } }} />
    </div>
  );
}

export default App;
