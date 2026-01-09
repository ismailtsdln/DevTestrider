import { ChevronRight, FileCode, CheckCircle2, XCircle, Clock } from 'lucide-react';

export function TestDetails() {
  // Mock data for now
  const packages = [
    { name: 'github.com/ismailtsdln/DevTestrider/internal/engine', status: 'PASS', duration: 0.45, tests: 23 },
    { name: 'github.com/ismailtsdln/DevTestrider/internal/server', status: 'PASS', duration: 0.12, tests: 8 },
    { name: 'github.com/ismailtsdln/DevTestrider/cmd', status: 'FAIL', duration: 0.05, tests: 2 },
  ];

  return (
    <div className="space-y-6">
      <div className="bg-slate-900/50 border border-slate-800 rounded-xl overflow-hidden backdrop-blur-sm">
        <div className="px-6 py-4 border-b border-slate-800 flex items-center justify-between">
            <h3 className="text-lg font-semibold text-slate-50">Test Packages</h3>
            <div className="text-sm text-slate-400">Total: 3 packages</div>
        </div>
        
        <div className="divide-y divide-slate-800">
            {packages.map((pkg, idx) => (
                <div key={idx} className="p-4 hover:bg-slate-800/30 transition-colors flex items-center justify-between group cursor-pointer">
                    <div className="flex items-center gap-4">
                        {pkg.status === 'PASS' 
                            ? <CheckCircle2 className="w-5 h-5 text-emerald-400" />
                            : <XCircle className="w-5 h-5 text-rose-400" />
                        }
                        <div>
                            <p className="font-medium text-slate-200 group-hover:text-indigo-300 transition-colors font-mono text-sm">{pkg.name}</p>
                            <div className="flex items-center gap-3 mt-1 text-xs text-slate-500">
                                <span className="flex items-center gap-1"><FileCode size={12}/> {pkg.tests} tests</span>
                                <span className="flex items-center gap-1"><Clock size={12}/> {pkg.duration}s</span>
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
