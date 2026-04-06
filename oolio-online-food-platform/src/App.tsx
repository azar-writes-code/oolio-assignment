import AppRouter from "@/routes";

function AppContent() {
  // handle auth listener
  // handle network listener
  // handle theme listener

  return (
    <>
      <AppRouter />
    </>
  );
}

function App() {
  return <AppContent />;
}

export default App;
