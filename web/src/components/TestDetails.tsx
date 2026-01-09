import { useEffect, useState } from 'react';
import { ChevronRight, FileCode, CheckCircle2, XCircle, Clock, Percent } from 'lucide-react';
import clsx from 'clsx';

export function TestDetails() {
  const [result, setResult] = useState<any>(null);

  const fetchLatest = async () => {
    try {
      const res = await fetch('/api/results/latest');
      if (res.ok) {
        const data = await res.json();
        setResult(data);
      }
    } catch (e) {
      console.error(e);
    }
  };

  useEffect(() => {
    fetchLatest();
    // Subscribe to SSE for updates
    const eventSource = new EventSource('/api/events');
    eventSource.onmessage = () => fetchLatest();
    return () => eventSource.close();
  }, []);

  if (!result) return <div className="p-8 text-center text-slate-500">Loading test detail...</div>;

  // Convert packages map to array
  const packages = Object.values(result.packages || {}).sort((a: any, b: any) => a.name.localeCompare(b.name));

  return (
    <div className="space-y-6">
      <div className="bg-slate-900/50 border border-slate-800 rounded-xl overflow-hidden backdrop-blur-sm">
        <div className="px-6 py-4 border-b border-slate-800 flex items-center justify-between">
            <h3 className="text-lg font-semibold text-slate-50">Test Packages</h3>
            <div className="text-sm text-slate-400">Total: {packages.length} packages</div>
        </div>
        
        <div className="divide-y divide-slate-800">
            {packages.map((pkg: any, idx: number) => (
                <div key={idx} className="p-4 hover:bg-slate-800/30 transition-colors flex items-center justify-between group cursor-pointer">
                    <div className="flex items-center gap-4">
                        {pkg.status === 'PASS' 
                            ? <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                            : <XCircle className="w-5 h-5 text-rose-400" />
                        }
                        <div>
                            <p className="font-medium text-slate-200 group-hover:text-indigo-300 transition-colors font-mono text-sm">{pkg.name}</p>
                            <div className="flex items-center gap-3 mt-1 text-xs text-slate-500">
                                <span className="flex items-center gap-1"><FileCode size={12}/> {pkg.tests?.length || 0} tests</span>
                                <span className="flex items-center gap-1"><Clock size={12}/> {pkg.duration.toFixed(3)}s</span>
                                {pkg.coverage > 0 && (
                                    <span className={clsx("flex items-center gap-1 font-semibold", 
                                        pkg.coverage > 80 ? "text-emerald-400" : 
                                        pkg.coverage > 50 ? "text-amber-400" : "text-rose-400"
                                    )}>
                                        <Percent size={12}/> {pkg.coverage.toFixed(1)}%
                                    </span>
                                )}
                            </div>
                        </div>
                    </div>
                    <ChevronRight className="w-5 h-5 text-slate-600 group-hover:text-slate-400" />
                </div>
            ))}
        </div>
      </div>
    </div>
  );
}
