import React from 'react';
import { useForm } from 'react-hook-form';
import { zodResolver } from '@hookform/resolvers/zod';
import * as z from 'zod';
import { motion } from 'framer-motion';
import { Save, AlertCircle } from 'lucide-react';

const policySchema = z.object({
  name: z.string().min(3, { message: 'Policy name must be at least 3 characters' }),
  method: z.enum(['tools/list', 'tools/call', 'resources/list', 'resources/read', 'prompts/list', 'prompts/get']),
  role: z.string().min(1, { message: 'Role is required' }),
  description: z.string().optional(),
});

const PolicyEditor = () => {
  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm({
    resolver: zodResolver(policySchema),
    defaultValues: {
      method: 'tools/call',
      role: 'developer',
    },
  });

  const onSubmit = (data) => {
    console.log('Policy submitted:', data);
    alert('Policy draft saved! (Simulation)');
  };

  return (
    <motion.div
      initial={{ opacity: 0, scale: 0.95 }}
      animate={{ opacity: 1, scale: 1 }}
      className="card max-w-2xl mx-auto"
    >
      <div className="flex items-center gap-3 mb-lg">
        <div className="bg-primary/10 p-2 rounded-md">
          <Save className="text-primary w-6 h-6" />
        </div>
        <h2 className="text-2xl font-bold">New Security Policy</h2>
      </div>

      <form onSubmit={handleSubmit(onSubmit)} className="space-y-md">
        <div>
          <label className="block text-sm font-bold text-tertiary mb-2 uppercase tracking-wider">
            Policy Name
          </label>
          <input
            {...register('name')}
            className={`input w-full ${errors.name ? 'border-error ring-1 ring-error' : ''}`}
            placeholder="e.g., Allow File Access"
          />
          {errors.name && (
            <p className="mt-1 text-sm text-error flex items-center gap-1">
              <AlertCircle className="w-3 h-3" /> {errors.name.message}
            </p>
          )}
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-md">
          <div>
            <label className="block text-sm font-bold text-tertiary mb-2 uppercase tracking-wider">
              MCP Method
            </label>
            <select
              {...register('method')}
              className="input w-full"
            >
              <option value="tools/list">tools/list</option>
              <option value="tools/call">tools/call</option>
              <option value="resources/list">resources/list</option>
              <option value="resources/read">resources/read</option>
              <option value="prompts/list">prompts/list</option>
              <option value="prompts/get">prompts/get</option>
            </select>
          </div>

          <div>
            <label className="block text-sm font-bold text-tertiary mb-2 uppercase tracking-wider">
              Target Role
            </label>
            <input
              {...register('role')}
              className={`input w-full ${errors.role ? 'border-error ring-1 ring-error' : ''}`}
              placeholder="e.g., admin"
            />
            {errors.role && (
              <p className="mt-1 text-sm text-error flex items-center gap-1">
                <AlertCircle className="w-3 h-3" /> {errors.role.message}
              </p>
            )}
          </div>
        </div>

        <div>
          <label className="block text-sm font-bold text-tertiary mb-2 uppercase tracking-wider">
            Description
          </label>
          <textarea
            {...register('description')}
            className="input w-full h-24"
            placeholder="Explain the purpose of this policy..."
          />
        </div>

        <div className="pt-sm border-t border-border/50 flex justify-end gap-md">
          <button type="button" className="text-sm font-bold text-tertiary hover:underline">
            Cancel
          </button>
          <button type="submit" className="btn-primary !w-auto !px-lg">
            Create Policy
          </button>
        </div>
      </form>
    </motion.div>
  );
};

export default PolicyEditor;
