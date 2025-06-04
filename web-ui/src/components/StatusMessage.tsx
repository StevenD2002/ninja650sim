import { useAppSelector } from "../hooks/redux"
export const StatusMessage = () => {
  const ui = useAppSelector((state: any) => state.motorcycle.ui)
  return (
    <>
     {ui.statusMessage && Date.now() - ui.statusMessageTime < 3000 && (
        <div style={{
          position: 'fixed',
          bottom: '20px',
          left: '50%',
          transform: 'translateX(-50%)',
          padding: '12px 24px',
          borderRadius: '8px',
          fontWeight: '600',
          backgroundColor: '#374151',
          border: '1px solid #4b5563',
          color: 'white',
          zIndex: 1000
        }}>
          {ui.statusMessage}
        </div>
      )}
    </>
  )
}
