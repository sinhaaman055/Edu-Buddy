A real-time, collaborative co-op testing platform modeled after Testbook and Discord, designed for students to take competitive mock exams synchronously with their peers.
EduBuddy eliminates the isolation of standard online mock tests by introducing a shared, synchronized testing room. Peers can join a room, select an exam, and instantly face the exact same question set with a globally synchronized server-side timer.
While the test is active, each student’s environment remains completely isolated—individual answers and inputs are tracked silently on the backend without affecting the other peer's screen, preserving pure competitive integrity. Once the test wraps up or the timer hits zero, the platform triggers a synchronous "Big Reveal", rendering a side-by-side comparison matrix showing exact accuracy metrics, question-specific time consumption, and immediate solutions.

🌟 Key Architectural Features
Synchronized Testing Lobbies: Dynamic room creation and real-time state synchronization via WebSockets, maintaining uniform test states across multiple concurrent clients.

Host Control Center: Room owners have full administrative privileges to customize the environment, including mid-test chat authorization to prevent collusion.

Dual Question Pipelines:

Agentic AI Engine: An autonomous API pipeline that searches, curates, and converts raw exam patterns/PYQs into structured JSON data on the fly.

Coaching PDF Parser: An upload-and-parse system allowing room hosts to transform traditional mock test PDFs into fully interactive digital exams.

Granular Response Tracking: Event-driven backend system that processes individual user choices in memory without broadcasting them to room peers until the final evaluation phase.
