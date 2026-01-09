import { useEffect, useState } from 'react';
import { ChevronRight, FileCode, CheckCircle2, XCircle, Clock, Percent } from 'lucide-react';
import clsx from 'clsx';

interface TestCase {
  name: string;
  duration: number;
  status: string;
}

interface PackageResult {
  name: string;
  duration: number;
  status: string;
  tests: TestCase[];
  coverage: number;
}

interface TestResult {
  packages: Record<string, PackageResult>;
  issues?: string[];
}

export function TestDetails() {
  const [result, setResult] = useState<TestResult | null>(null);
  const [expanded, setExpanded] = useState<Record<string, boolean>>({});

  useEffect(() => {
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

    fetchLatest();
    // Subscribe to SSE for updates
    const eventSource = new EventSource('/api/events');
    eventSource.onmessage = () => fetchLatest();
    return () => eventSource.close();
  }, []);

  const toggleExpand = (name: string) => {
    setExpanded(prev => ({...prev, [name]: !prev[name]}));
  };

  if (!result) return <div className="p-8 text-center text-slate-500">Loading test detail...</div>;

  // Convert packages map to array
  const packages = Object.values(result.packages || {}).sort((a: PackageResult, b: PackageResult) => a.name.localeCompare(b.name));

  return (
    <div className="space-y-6">
      
      {/* Analysis Issues */}
      {result.issues && result.issues.length > 0 && (
          <div className="bg-amber-950/20 border border-amber-900/50 rounded-xl overflow-hidden backdrop-blur-sm">
             <div className="px-6 py-4 border-b border-amber-900/50 flex items-center gap-2">
                 <XCircle className="w-5 h-5 text-amber-500" />
                 <h3 className="text-lg font-semibold text-amber-400">Analysis Issues</h3>
             </div>
             <div className="p-4 space-y-2">
                 {result.issues.map((issue, idx) => (
                     <div key={idx} className="font-mono text-sm text-amber-200/80 bg-amber-900/20 p-2 rounded border border-amber-900/30">
                         {issue}
                     </div>
                 ))}
             </div>
          </div>
      )}

      <div className="bg-slate-900/50 border border-slate-800 rounded-xl overflow-hidden backdrop-blur-sm">
        <div className="px-6 py-4 border-b border-slate-800 flex items-center justify-between">
            <h3 className="text-lg font-semibold text-slate-50">Test Packages</h3>
            <div className="text-sm text-slate-400">Total: {packages.length} packages</div>
        </div>
        
        <div className="divide-y divide-slate-800">
            {packages.map((pkg: PackageResult, idx: number) => (
                <div key={idx} className="group">
                    <div 
                        className="p-4 hover:bg-slate-800/30 transition-colors flex items-center justify-between cursor-pointer"
                        onClick={() => toggleExpand(pkg.name)}
                    >
                        <div className="flex items-center gap-4">
                            <ChevronRight className={clsx("w-5 h-5 text-slate-500 transition-transform", expanded[pkg.name] ? "rotate-90" : "")} />
                            {pkg.status === 'PASS' 
                                ? <CheckCircle2 className="w-5 h-5 text-emerald-500" />
                                : <XCircle className="w-5 h-5 text-rose-500" />
                            }
                            <div>
                                <h4 className="font-medium text-slate-200">{pkg.name}</h4>
                                <div className="flex items-center gap-3 text-xs text-slate-500 mt-1">
                                    <span className="flex items-center gap-1"><Clock className="w-3 h-3" /> {(pkg.duration).toFixed(2)}s</span>
                                    <span className="flex items-center gap-1"><FileCode className="w-3 h-3" /> {pkg.tests?.length || 0} tests</span>
                                </div>
                            </div>
                        </div>
                        <div className="flex items-center gap-4">
                            {pkg.coverage > 0 && (
                                <div className={clsx(
                                    "flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-medium border",
                                    pkg.coverage >= 80 ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/20" :
                                    pkg.coverage >= 50 ? "bg-amber-500/10 text-amber-400 border-amber-500/20" :
                                    "bg-rose-500/10 text-rose-400 border-rose-500/20"
                                )}>
                                    <Percent className="w-3 h-3" />
                                    {pkg.coverage.toFixed(1)}%
                                </div>
                            )}
                        </div>
                    </div>
                    
                    {/* Expanded Test Cases */}
                    {expanded[pkg.name] && pkg.tests && pkg.tests.length > 0 && (
                        <div className="bg-slate-950/30 px-4 py-2 border-t border-slate-800/50">
                            {pkg.tests.map((test, tIdx) => (
                                <div key={tIdx} className="flex items-center justify-between py-2 border-b border-slate-800/30 last:border-0 pl-9">
                                    <div className="flex items-center gap-3">
                                         {test.status === 'PASS' 
                                            ? <div className="w-2 h-2 rounded-full bg-emerald-500" />
                                            : <div className="w-2 h-2 rounded-full bg-rose-500" />
                                         }
                                         <span className={clsx("text-sm font-mono", test.status === 'PASS' ? "text-slate-400" : "text-rose-300")}>
                                            {test.name}
                                         </span>
                                    </div>
                                    <span className="text-xs text-slate-500 font-mono">{(test.duration).toFixed(3)}s</span>
                                </div>
                            ))}
                        </div>
                    )}
                </div>
            ))}
        </div>
      </div>
    </div>
  );
}
