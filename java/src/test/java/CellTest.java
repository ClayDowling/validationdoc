import org.junit.jupiter.api.Test;
import org.junit.jupiter.params.ParameterizedTest;
import org.junit.jupiter.params.provider.ValueSource;

import static org.assertj.core.api.Assertions.assertThat;
import static org.junit.jupiter.api.Assertions.assertEquals;

public class CellTest {

    @Test
    public void Cell_ByDefault_SetsStateFromConstructor() {
        Cell c = new Cell(true);
        assertEquals(true, c.isAlive());

        Cell d = new Cell(false);
        assertThat(d.isAlive()).isFalse();
    }

    // Requirement CELL-1
    @ParameterizedTest
    @ValueSource(ints = { 0, 1 })
    void Cell_LiveWithFewerThanTwoNeighbors_Dies(int neighbor) {
        Cell c = new Cell(true);
        assertThat(c.Lives(neighbor)).isFalse();
    }

    // Requirement CELL-2
    @ParameterizedTest
    @ValueSource(ints = { 2, 3 })
    void LiveCell_WithTwoOrThreeNeighbors_Lives(int neighbors) {
        Cell c = new Cell(true);
        assertThat(c.Lives(neighbors)).isTrue();
    }

    @Test
    // Requirement CELL-3
    void DeadCell_WithExactly3Neighbors_Lives() {
        Cell c = new Cell(false);
        assertThat(c.Lives(3)).isTrue();
    }

    @ParameterizedTest
    @ValueSource(ints = { 4, 5, 6, 7, 8 })
    // Requirement CELL-4
    void LiveCell_With4OrMoreNeighbors_Dies(int neighbors) {
        Cell c = new Cell(true);
        assertThat(c.Lives(neighbors)).isFalse();
    }

    // Requirement CELL-5
    @ParameterizedTest
    @ValueSource(ints = { 0, 1, 2, 4, 5, 6, 7, 8 })
    void DeadCell_WithNot3Neighbors_StaysDead(int neighbors) {
        Cell c = new Cell(false);
        assertThat(c.Lives(neighbors)).isFalse();
    }

}
