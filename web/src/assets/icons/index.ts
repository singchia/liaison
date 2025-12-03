const requireAll = (requireContext: any) => {
  return requireContext.keys().map(requireContext);
};

const req = (require as any).context('@/assets/icons', true, /\.svg$/);

requireAll(req);
