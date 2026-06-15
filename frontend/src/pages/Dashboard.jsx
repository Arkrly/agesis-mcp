import React, { useState, useEffect } from 'react';
import { motion } from 'framer-motion';
import { ShieldCheck, ShieldAlert, Activity, Cpu, FileText, CheckCircle, XCircle } from 'lucide-react';
import api from '../api/client';

const Dashboard = () => {
  const [summary, setSummary] = useState(null);
  const [auditLogs, setAuditLogs] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchData = async () => {
      try {
        const [summaryRes, auditRes] = await Promise.all([
          api.get('/api/summary'),
          api.get('/api/audit')
        ]);
        setSummary(summaryRes.data);
        setAuditLogs(auditRes.data || []);
      } catch (err) {
        console.error('Failed to fetch dashboard data', err);
      } finally {
        setLoading(false);
      }
    };

    fetchData();
  }, []);

  if (loading) {
    return (
      <div className="flex items-center justify-center h-64">
        <div className="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-primary"></div>
      </div>
    );
  }

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: {
        staggerChildren: 0.1
      }
    }
  };

  const itemVariants = {
    hidden: { y: 20, opacity: 0 },
    visible: { y: 0, opacity: 1 }
  };

  return (
    <motion.div
      variants={containerVariants}
      initial="hidden"
      animate="visible"
      className="space-y-xl"
    >
      <section>
        <div className="flex items-center justify-between mb-md">
          <h1 className="text-3xl font-bold">System Overview</h1>
          <div className="flex items-center gap-2">
            <div className={`h-3 w-3 rounded-full ${summary?.status === 'ok' ? 'bg-green-500' : 'bg-red-500'} animate-pulse`} />
            <span className="text-sm font-semibold uppercase tracking-wider text-tertiary">
              Status: {summary?.status || 'Unknown'}
            </span>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-3 gap-md">
          <motion.div variants={itemVariants} className="card flex flex-col items-center text-center">
            <Activity className="w-10 h-10 text-primary mb-sm" />
            <h3 className="text-lg font-bold">Proxy Status</h3>
            <p className="text-sm text-tertiary mt-2">Active and Intercepting</p>
            <div className="chip mt-md">Healthy</div>
          </motion.div>

          <motion.div variants={itemVariants} className="card flex flex-col items-center text-center">
            <ShieldCheck className="w-10 h-10 text-primary mb-sm" />
            <h3 className="text-lg font-bold">Policies</h3>
            <p className="text-sm text-tertiary mt-2">Enforcing Zero-Trust</p>
            <div className="chip mt-md">8 Rules Active</div>
          </motion.div>

          <motion.div variants={itemVariants} className="card flex flex-col items-center text-center">
            <Cpu className="w-10 h-10 text-primary mb-sm" />
            <h3 className="text-lg font-bold">Inference</h3>
            <p className="text-sm text-tertiary mt-2">Semantic Inspection Engine</p>
            <div className="chip mt-md">v0.1.0-alpha</div>
          </motion.div>
        </div>
      </section>

      <section>
        <div className="flex items-center justify-between mb-md">
          <h2 className="text-2xl font-bold">Recent Security Events</h2>
          <button className="text-sm font-bold text-primary hover:underline flex items-center gap-1">
            View All Logs <FileText className="w-4 h-4" />
          </button>
        </div>

        <motion.div variants={itemVariants} className="card overflow-hidden !p-0">
          <div className="overflow-x-auto">
            <table className="w-full text-left border-collapse">
              <thead className="bg-surface-muted border-b border-border/50">
                <tr>
                  <th className="px-md py-sm text-xs font-bold uppercase text-tertiary">Timestamp</th>
                  <th className="px-md py-sm text-xs font-bold uppercase text-tertiary">Agent ID</th>
                  <th className="px-md py-sm text-xs font-bold uppercase text-tertiary">Method</th>
                  <th className="px-md py-sm text-xs font-bold uppercase text-tertiary">Decision</th>
                  <th className="px-md py-sm text-xs font-bold uppercase text-tertiary">Reason</th>
                </tr>
              </thead>
              <tbody className="divide-y divide-border/30">
                {auditLogs.length > 0 ? (
                  auditLogs.map((log, index) => (
                    <tr key={index} className="hover:bg-surface-muted/50 transition-colors">
                      <td className="px-md py-md text-sm">
                        {new Date(log.ts).toLocaleString()}
                      </td>
                      <td className="px-md py-md text-sm font-mono text-primary font-semibold">
                        {log.agent_id}
                      </td>
                      <td className="px-md py-md text-sm">
                        <span className="bg-neutral px-2 py-1 rounded text-xs font-bold border border-border/50">
                          {log.method}
                        </span>
                      </td>
                      <td className="px-md py-md text-sm">
                        <div className="flex items-center gap-1">
                          {log.allowed ? (
                            <>
                              <CheckCircle className="w-4 h-4 text-green-500" />
                              <span className="text-green-600 font-bold">ALLOWED</span>
                            </>
                          ) : (
                            <>
                              <ShieldAlert className="w-4 h-4 text-red-500" />
                              <span className="text-red-600 font-bold">BLOCKED</span>
                            </>
                          )}
                        </div>
                      </td>
                      <td className="px-md py-md text-sm italic text-tertiary">
                        {log.reason || 'Allowed by policy'}
                      </td>
                    </tr>
                  ))
                ) : (
                  <tr>
                    <td colSpan="5" className="px-md py-xl text-center text-tertiary">
                      No security events recorded yet.
                    </td>
                  </tr>
                )}
              </tbody>
            </table>
          </div>
        </motion.div>
      </section>
    </motion.div>
  );
};

export default Dashboard;
