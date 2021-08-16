using FluentAssertions;
using gameoflife;
using Xunit;

namespace gameoflife_tests
{
    public class BoardTests
    {
        private int initializerCalled;

        public BoardTests()
        {
            initializerCalled = 0;
        }
        
        [Fact]
        // Requirement GAME-1
        public void Board_ByDefault_RandomlyInitializesState()
        {
            var board = new Board(3, 3, MockInitializer);
            initializerCalled.Should().Be(9);
        }

        public bool MockInitializer()
        {
            initializerCalled++;
            return true;
        }
    }
}