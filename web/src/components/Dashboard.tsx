import { useEffect, useState } from 'react';
import { AreaChart, Area, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer } from 'recharts';
import { PlayCircle, XCircle, CheckCircle2, Activity } from 'lucide-react';
import clsx from 'clsx';

// Type definitions matching backend
interface TestResult {
  timestamp: string;
  total_tests: number;
  passed_tests: number;
  failed_tests: number;
  skipped_tests: number;
  duration: number;
  success: boolean;
}

const StatCard = ({ title, value, sub, icon: Icon, color }: any) => (
  <div className="bg-slate-900/50 border border-slate-800 p-6 rounded-xl backdrop-blur-sm">
    <div className="flex items-start justify-between">
      <div>
        <p className="text-sm font-medium text-slate-400">{title}</p>
        <h3 className="text-3xl font-bold text-slate-50 mt-2">{value}</h3>
        <p className={clsx("text-xs mt-1 font-mono", color)}>{sub}</p>
      </div>
      <div className={clsx("p-3 rounded-lg bg-opacity-10", color.replace('text-', 'bg-'))}>
        <Icon className={clsx("w-6 h-6", color)} />
      </div>
    </div>
  </div>
);

export function Dashboard() {
  const [data, setData] = useState<TestResult | null>(null);
  const [history, setHistory] = useState<any[]>([]);

  useEffect(() => {
     const fetchLatest = async () => {
      try {
        const res = await fetch('/api/results/latest');
        if (res.ok) {
          const result = await res.json();
          if (result) {
              setData(result);
              setHistory(prev => [...prev.slice(-19), {
                  name: new Date(result.timestamp).toLocaleTimeString(),
                  pass: result.passed_tests,
                  fail: result.failed_tests
              }]);
          }
        }
      } catch (e) {
        console.error(e);
      }
    };

    // Setup SSE connection
    const eventSource = new EventSource('/api/events');
    
    eventSource.onmessage = (event) => {
      console.log("Event:", event.data);
      fetchLatest();
    };

    fetchLatest();

    return () => {
      eventSource.close();
    };
  }, []);

  return (
    <div className="space-y-6">
      {/* Stats Grid */}
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        <StatCard 
          title="Total Tests" 
          value={data?.total_tests || 0} 
          sub="Across 12 packages" 
          icon={PlayCircle} 
          color="text-indigo-400" 
        />
        <StatCard 
          title="Passed" 
          value={data?.passed_tests || 0} 
          sub="98% Success Rate" 
          icon={CheckCircle2} 
          color="text-emerald-400" 
        />
        <StatCard 
          title="Failed" 
          value={data?.failed_tests || 0} 
          sub="Requires Attention" 
          icon={XCircle} 
          color="text-rose-400" 
        />
        <StatCard 
          title="Duration" 
          value={`${(data?.duration || 0).toFixed(2)}s`} 
          sub="Average: 1.2s" 
          icon={Activity}
          color="text-amber-400" 
        />
      </div>

      {/* Main Chart Section */}
      <div className="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <div className="lg:col-span-2 bg-slate-900/50 border border-slate-800 rounded-xl p-6 backdrop-blur-sm">
          <h3 className="text-lg font-semibold text-slate-50 mb-6">Execution Trend</h3>
          <div className="h-80">
            <ResponsiveContainer width="100%" height="100%">
              <AreaChart data={history.length ? history : [{name:'0', pass:0, fail:0}]}>
                <defs>
                  <linearGradient id="colorPass" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#34d399" stopOpacity={0.1}/>
                    <stop offset="95%" stopColor="#34d399" stopOpacity={0}/>
                  </linearGradient>
                  <linearGradient id="colorFail" x1="0" y1="0" x2="0" y2="1">
                    <stop offset="5%" stopColor="#f43f5e" stopOpacity={0.1}/>
                    <stop offset="95%" stopColor="#f43f5e" stopOpacity={0}/>
                  </linearGradient>
                </defs>
                <CartesianGrid strokeDasharray="3 3" stroke="#1e293b" vertical={false} />
                <XAxis dataKey="name" stroke="#64748b" fontSize={12} tickLine={false} axisLine={false} />
                <YAxis stroke="#64748b" fontSize={12} tickLine={false} axisLine={false} />
                <Tooltip 
                    contentStyle={{ backgroundColor: '#0f172a', borderColor: '#1e293b', color: '#f8fafc' }} 
                    itemStyle={{ color: '#f8fafc' }}
                />
                <Area type="monotone" dataKey="pass" stroke="#34d399" strokeWidth={2} fillOpacity={1} fill="url(#colorPass)" />
                <Area type="monotone" dataKey="fail" stroke="#f43f5e" strokeWidth={2} fillOpacity={1} fill="url(#colorFail)" />
              </AreaChart>
            </ResponsiveContainer>
          </div>
        </div>

        {/* Recent Activity / Logs */}
        <div className="bg-slate-900/50 border border-slate-800 rounded-xl p-6 backdrop-blur-sm flex flex-col">
            <h3 className="text-lg font-semibold text-slate-50 mb-4">Recent Activity</h3>
            <div className="flex-1 overflow-y-auto space-y-4 pr-2 custom-scrollbar">
                {!data && <p className="text-slate-500 text-sm">No tests run yet. Waiting for file changes...</p>}
                {data && (
                    <div className="space-y-3">
                         <div className="flex items-center gap-3 p-3 rounded-lg bg-slate-800/40 border border-slate-800">
                            {data.success ? <CheckCircle2 className="text-emerald-400 w-5 h-5"/> : <XCircle className="text-rose-400 w-5 h-5"/>}
                            <div>
                                <p className="text-sm font-medium text-slate-200">{data.success ? 'Test Suite Passed' : 'Test Suite Failed'}</p>
                                <p className="text-xs text-slate-500">{new Date(data.timestamp).toLocaleTimeString()}</p>
                            </div>
                        </div>
                    </div>
                )}
            </div>
        </div>
      </div>
    </div>
  );
}
