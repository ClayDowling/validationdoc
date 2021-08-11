using FluentAssertions;
using gameoflife;
using Xunit;

namespace gameoflife_tests
{
    public class CellTests
    {
        [Theory]
        [InlineData(0)]
        [InlineData(1)]
        /// Requirement CELL-1
        public void Lives_GivenZeroNeighbors_ReturnsFalse(int neighbors)
        {
            Cell c = new Cell(true);
            c.Lives(neighbors).Should().BeFalse();
        }

        [Theory]
        [InlineData(2)]
        [InlineData(3)]
        /// Requirement CELL-2
        public void Lives_GiveLiveCellWith2Or3Neighbors_ReturnsTrue(int neighbors)
        {
            Cell c = new Cell(true);
            c.Lives(neighbors).Should().BeTrue();
        }

        [Fact]
        /// Requirement CELL-3
        public void Lives_GivenDeadCellWith3Neighbors_ReturnsTrue() {
            Cell c = new Cell(false);
            c.Lives(3).Should().BeTrue();
        }

        [Theory]
        [InlineData(4)]
        [InlineData(5)]
        [InlineData(6)]
        [InlineData(7)]
        [InlineData(8)]
        /// Requirement CELL-4
        public void Lives_GivenLiveCellWith4OrMoveNeighbors_ReturnsFalse(int neighbors) {
            Cell c = new Cell(true);
            c.Lives(neighbors).Should().BeFalse();
        }

        /// Requirement CELL-5
        [Theory]
        [InlineData(0)]
        [InlineData(1)]
        [InlineData(2)]
        [InlineData(4)]
        [InlineData(5)]
        [InlineData(6)]
        [InlineData(7)]
        [InlineData(8)]
        public void Lives_GivenDeadCellWithNot3Neighbors_ReturnsFalse(int neighbors)
        {
            Cell c = new Cell(false);
            c.Lives(neighbors).Should().BeFalse();
        }

    }
}