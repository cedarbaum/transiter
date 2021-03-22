"""Add transfers table

Revision ID: 65551af9cd46
Revises: 423495a41dd6
Create Date: 2020-06-04 10:56:09.945967

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = "65551af9cd46"
down_revision = "423495a41dd6"
branch_labels = None
depends_on = None


def upgrade():
    # ### commands auto generated by Alembic - please adjust! ###
    op.create_table(
        "transfer",
        sa.Column("pk", sa.Integer(), nullable=False),
        sa.Column("source_pk", sa.Integer(), nullable=True),
        sa.Column("system_pk", sa.Integer(), nullable=True),
        sa.Column("from_stop_pk", sa.Integer(), nullable=True),
        sa.Column("to_stop_pk", sa.Integer(), nullable=True),
        sa.Column(
            "type",
            sa.Enum(
                "RECOMMENDED",
                "COORDINATED",
                "POSSIBLE",
                "NO_TRANSFER",
                name="type",
                native_enum=False,
                create_constraint=True,
            ),
            nullable=False,
        ),
        sa.Column("min_transfer_time", sa.Integer(), nullable=True),
        sa.ForeignKeyConstraint(["from_stop_pk"], ["stop.pk"],),
        sa.ForeignKeyConstraint(["system_pk"], ["system.pk"],),
        sa.ForeignKeyConstraint(["source_pk"], ["feed_update.pk"],),
        sa.ForeignKeyConstraint(["to_stop_pk"], ["stop.pk"],),
        sa.PrimaryKeyConstraint("pk"),
    )
    op.create_index(
        op.f("ix_transfer_from_stop_pk"), "transfer", ["from_stop_pk"], unique=False
    )
    op.create_index(
        op.f("ix_transfer_system_pk"), "transfer", ["system_pk"], unique=False
    )
    op.create_index(
        op.f("ix_transfer_source_pk"), "transfer", ["source_pk"], unique=False
    )
    op.create_index(
        op.f("ix_transfer_to_stop_pk"), "transfer", ["to_stop_pk"], unique=False
    )

    op.drop_constraint("result", "feed_update")


def downgrade():
    pass
