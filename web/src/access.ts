export default (initialState: { currentUser?: API.CurrentUser }) => {
  const canSeeAdmin = !!(initialState && initialState.currentUser);
  return {
    canSeeAdmin,
  };
};
